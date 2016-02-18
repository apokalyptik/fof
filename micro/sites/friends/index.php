<?php

define( SLACK_TOKEN, 'xoxp-...' );
define( CACHE_KEY, "fof-friend-wall-" . filemtime( realpath( __FILE__ ) ) );

ob_start();
header( 'Expires: '.gmdate('D, d M Y H:i:s', time()+3600).'GMT');
if ( empty( $_GET['debug'] ) ) {
	$mc = new Memcached;
	$mc->addServer("127.0.0.1", 11211);
	if ( $cache = unserialize( $mc->get( CACHE_KEY ) ) ) {
		header( 'X-Cached: true' );
		header( 'X-Cache-Key: ' . filemtime( realpath( __FILE__ ) ) );
		header( 'Etag: ' . md5( $cache->out ));
		header( 'Last-Modified: '.gmdate( 'D, d M Y H:i:s', $out->when ) );
		echo $cache->out;
		return;
	}
}
?><html lang="en">
<head>
	<meta charset="utf-8">
	<meta http-equiv="X-UA-Compatible" content="IE=edge">
	<title>FoF Friends List Helper</title>
	<script src="//cdnjs.cloudflare.com/ajax/libs/jquery/2.1.3/jquery.min.js"></script>
	<style>
	img.picon {
		vertical-align: text-top;
		float: left;
		margin-right: 5px;
	}
	div.pcontainer {
		/* clear: both; */
		/*height: 55px;*/
		height: 48px;
		padding: 5px;
		margin: 5px;
		background: #f7f7f7;
		width: 24em;
		-moz-border-radius: 10px;
		-webkit-border-radius: 10px;
		border-radius: 10px; /* future proofing */
		-khtml-border-radius: 10px; /* for old Konqueror browsers */
		float: left;
	}
	div.pname {
		text-overflow: ellipsis;
		font-weight: bold;
		font-size: 1.25em;
	}
	div.pcontainer.done {
		background: #DDD;
		opacity: 0.25;
	}
	a {
		font-weight: bold;
		background-color: #107c10;
		color: white;
		text-decoration: none;
		border: 1px solid black;
		padding: 2px 10px 2px 10px;
		-moz-border-radius: 10px;
		-webkit-border-radius: 10px;
		border-radius: 10px; /* future proofing */
		-khtml-border-radius: 10px; /* for old Konqueror browsers */
	}
	@media only screen 
	and (min-device-width : 320px) 
	and (max-device-width : 736px) {
		div.pcontainer input {
			min-height: 48px;
			min-width: 48px;
		}
		div.pcontainer {
			width: 100%;
			font-size: 2.5em;
		}
	}
</style>
<script>
	function supports_html5_storage() {
		try {
			return 'localStorage' in window && window['localStorage'] !== null;
		} catch (e) {
			return false;
		}
	}
function didClick(id) {
		if ( !supports_html5_storage() ) {
			return false;
		}
		if ( $("#did-" + id).prop('checked') ) {
			localStorage["did-"+id] = true;
			$("div.did-"+id).addClass("done");
		} else {
			$("div.did-"+id).removeClass("done");
			localStorage.removeItem("did-"+id);
		}
	}
	$(window).ready(function() {
		var inputs = $('input.did');
		$.each(inputs, function(i, v) {
			if ( !supports_html5_storage() ) {
				$(v).delete();
			}
			var id = $(v).attr("id");
			var did = localStorage[id];
			if ( localStorage[id] ) {
				$("div."+id).addClass("done");
				$(v).prop('checked', did);
				$("div."+id).appendTo("#alldone");
				//$("div."+id+":first").remove();

			}
			$(v).prop('disabled', false);
		});
	});
</script>
</head>
<body>
<?php


$ch = curl_init();
curl_setopt( $ch, CURLOPT_URL, 'http://127.0.0.1:8890/seen.json' );
curl_setopt( $ch, CURLOPT_RETURNTRANSFER, true );
$res = curl_exec( $ch );
$status = curl_getinfo( $ch, CURLINFO_HTTP_CODE );
$seen = null;
if ( $status == 200 ) {
	$seen = json_decode( $res );
	if ( !$seen )
		$seen = null;
}

$ch = curl_init(); 
curl_setopt( $ch, CURLOPT_URL, 'https://slack.com/api/users.list' );
curl_setopt( $ch, CURLOPT_RETURNTRANSFER, true );
curl_setopt( $ch, CURLOPT_POST, true );
curl_setopt( $ch, CURLOPT_POSTFIELDS, array( 'token' => SLACK_TOKEN ));
$res = curl_exec( $ch );
$status = curl_getinfo( $ch, CURLINFO_HTTP_CODE ); 

if ( !$res )
	die( "Unable to get the user list from slack. Contact @demitriousk if the problem persists</body></html>" );

$data = json_decode( $res );

if ( !$data || empty( $data ) )
	die( "Unable to parse the user list from slack. Contact @demitriousk if the problem persists</body></html>" );

if ( !$data->ok )
	die( "Unable to parse the user list from slack. Contact @demitriousk if the problem persists</body></html>" );

if ( !empty( $_GET['debug'] ) ) {
	//echo '<pre>'.print_r( $data->members, true ).'</pre>';
}

foreach( $data->members as $idx => $member ) {
	if ( $seen && is_object( $seen ) ) {
	   if ( !property_exists( $seen, $member->id ) )
		   continue;
	   $seen->{$member->id}[10] = ' ';
	   $seen->{$member->id} = strtotime( substr( $seen->{$member->id}, 0, 19 ) );
	   if ( ( time() - ( 31 * 86400 ) ) > $seen->{$member->id} )
		   continue;
	}
	if ( !empty( $member->is_bot ) )
		continue;
	if ( !empty( $member->deleted ) )
		continue;
	printf(
		'<div class="pcontainer did-%s">
		<div class="pname">
		<input disabled class="did" type="checkbox" id="did-%s" onClick="return didClick(\'%s\');"/>%s</div>
		XBox Live link: <a onClick="$(\'#did-%s\').click();return true;" target="_blank" href="https://account.xbox.com/en-us/profile?gamerTag=%s">%s</a></div>', 
		$member->id,
		$member->id,
		$member->id,
		htmlentities( 
			strlen( $member->profile->real_name_normalized ) > 30 ? substr( $member->profile->real_name_normalized, 0,27 ) . "..." : $member->profile->real_name_normalized  ),
		$member->id,
		rawurlencode( $member->profile->first_name ),
		htmlentities( $member->profile->first_name ) );
}
?>
	<div id="alldone"></div></body></html><?php
if ( empty( $_GET['debug'] ) ) {
	$when = time();
	$out = ob_get_clean();
	header( 'Last-Modified: '.gmdate( 'D, d M Y H:i:s', $when ) );
	header( 'X-Cached: false' );
	header( 'X-Cache-Key: ' . filemtime( realpath( __FILE__ ) ) );
	header( 'Etag: ' . md5( $out ));
	$mc->set( CACHE_KEY, serialize( (object)array( 'out' => $out, 'when' => $when ) ), $when + 300 );
    header('Etag: ' . md5( $out ));
}
echo $out;
?>
