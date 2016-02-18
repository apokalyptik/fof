<?php

define( SLACK_WEBHOOK, 'https://hooks.slack.com/services/...' );

if ( strlen( $_POST['text'] ) > 40 ) {
	$_POST['text'] = substr( $_POST['text'], 0, 37 ) . '...';
}
$gif_url = sprintf(
	"http://%s/360/%s",
	$_SERVER['HTTP_HOST'],
	rawurlencode( $_POST['text'] )
);
$ch = curl_init( SLACK_WEBHOOK );
curl_setopt( $ch, CURLOPT_CUSTOMREQUEST, "POST" );
curl_setopt( $ch, CURLOPT_POSTFIELDS, array( 
	'payload' => json_encode( 
		array(
			'channel' => "#".$_POST['channel_name'],
			'text' => sprintf( "%s %s\r\n<%s>", $_POST['command'], $_POST['text'], $gif_url ),
			'username' => $_POST['user_name'],
			'icon_emoji' => ':thumbsup:',
			"unfurl_links" => true
		)
	)
) );
$res = curl_exec( $ch );
if ( ! $res ) {
	header("HTTP/1.0 500 Internal Error");
	die( "Post-Back to slack failed!" );
}
