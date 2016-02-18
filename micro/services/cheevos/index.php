<?php

$requested_text = substr( rawurldecode( $_SERVER['REQUEST_URI'] ), 5, 40 );

$hash = sha1( $requested_text );

$cache_dir = sprintf(
	"%s/cache/%s/%s",
	dirname( __FILE__ ),
	substr( $hash, 0, 1 ),
	substr( $hash, 2, 1 )
);

if ( !is_dir( $cache_dir ) ) {
	mkdir( $cache_dir, 0777, true );
}

$cache_file = sprintf(
	"%s/%s.png",
	$cache_dir,
	substr( $hash, 3 )
);

if ( file_exists( $cache_file ) ) {
	header( 'Content-type: image/png' );
	$fp = fopen( $cache_file, 'r' );
	fpassthru( $fp );
	die();
}

// Never Expire
header( "Cache-Control: no-cache, must-revalidate" );
header( "Expires: Sat, 26 Jul 1997 05:00:00 GMT" );

// Read the background image into memory
$imagick = new Imagick();
$imagick->readImage( dirname( __FILE__ ) . '/achievement.gif' );

// Setup our annotation settings
$draw = new ImagickDraw();
$draw->setFillColor( 'white' );
$draw->setFont('./ConvectionRegular.ttf');
$draw->setFontSize( 20 );

// Annotate the user input
$imagick->annotateImage( $draw, 75, 55, 0, $requested_text );

// Annotate the header
$draw->setFontSize( 25 );
$imagick->annotateImage( $draw, 75, 30, 0, "ACHIEVEMENT UNLOCKED" );

// Output the image
$imagick->setImageFormat('png');
header('Content-type: image/png');
file_put_contents( $cache_file, $imagick->getimageblob() );
$fp = fopen( $cache_file, 'r' );
fpassthru( $fp );
