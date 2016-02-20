<?php

define( SLACK_TOKEN, 'xoxp-...' );

ob_start(); 
?><html lang="en">
<head>
	<title>FoF Friends List Helper</title>
	<script src="//cdnjs.cloudflare.com/ajax/libs/jquery/2.1.3/jquery.min.js"></script>
	<style>
	img.picon {
		vertical-align: text-middle;
		float: left;
		margin-right: 5px;
		height: 2em;
		-moz-border-radius: 10px;
		-webkit-border-radius: 10px;
		border-radius: 10px; /* future proofing */
		-khtml-border-radius: 10px; /* for old Konqueror browsers */
	}
	div.pcontainer {
		/* clear: both; */
		/*height: 55px;*/
		padding: 5px;
		margin: 5px;
		background: #f7f7f7;
		width: 18em;
		-moz-border-radius: 10px;
		-webkit-border-radius: 10px;
		border-radius: 10px; /* future proofing */
		-khtml-border-radius: 10px; /* for old Konqueror browsers */
		float: left;
	}
	div.pname {
		font-weight: bold;
		font-size: 1.25em;
		padding-top: 4px;
	}
	div.pcontainer.done {
		background: #DDD;
		opacity: 0.25;
	}
	a {
		font-weight: bold;
		text-transform: capitalize;
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
		$('a').click(function() {
			$(this).parent().find('input').click(); // .prop('checked', true);
		});
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

$cache_key = "fof-friend-wall-" . $_SERVER['HTTP_HOST'];

if ( empty( $_GET['debug'] ) ) {
	$mc = new Memcached;
	$mc->addServer("127.0.0.1", 11211);
	if ( false || $out = $mc->get( $cache_key ) ) {
	header('X-Cached: true');
		die( $out );
	}
}

$ch = curl_init(); 
curl_setopt( $ch, CURLOPT_URL, 'https://slack.com/api/users.list' );
curl_setopt( $ch, CURLOPT_RETURNTRANSFER, true );
curl_setopt( $ch, CURLOPT_POST, true );
curl_setopt( $ch, CURLOPT_POSTFIELDS, array( 'token' => SLACK_TOKEN ) );
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
	if ( !empty( $member->is_bot ) )
		continue;
	if ( !empty( $member->deleted ) )
		continue;
	if ( $member->id == "U053GJH59" ) // teambot
		continue;
	printf(
		'<div class="pcontainer did-%s">
		<img src="%s" class="picon">
		<div class="pname">
		<input disabled class="did" type="checkbox" id="did-%s" onClick="return didClick(\'%s\');"/>
		<!-- %s -->
		<!-- XBox Live link: -->
		<a target="_blank" href="https://account.xbox.com/en-us/profile?gamerTag=%s">%s</a>
		</div></div>', 
		$member->id,
		$member->profile->image_48, 
		$member->id,
		$member->id,
		htmlentities( $member->profile->first_name ),
		rawurlencode( $member->profile->first_name ),
		htmlentities( $member->profile->first_name ) );
}
?>
	<div id="alldone"></div></body></html><?php
if ( empty( $_GET['debug'] ) ) {
	$out = ob_get_clean(); 
	$mc->set( $cache_key, $out, time() + 60 );
}
die($out);
?>
