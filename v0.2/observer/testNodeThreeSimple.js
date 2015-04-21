// fs = require("fs")

THREE = require("./")
join = require("path").join

var camera, scene, renderer;
var frameNumber = 0;
var simTime = 0.0;
var deltaTime = 0.005;
var ground;
var velocity, angularVelocity; 
var pendulum, pendulum2;
var mass = 1.0;
var gravity = 10.0;
var theta = (Math.PI*2.0) - Math.PI/4.0;
var length = 1.0;
var b_height = 2*length, b_depth = 0.15, b_width = 0.15;
// var axis = Vector.create([0, 0.0, 1.0]);
// var position = Vector.create([0, -1.0, 0]);
var position = ([Math.cos(theta)*length, Math.sin(theta)*length, 0]);

var drawAngularForce, drawGavityForce, drawVelocity, drawVelocity, drawAngularVelocity, drawAcceleration, drawAngularAcceleration;
var drawTorqueForce;
var oldAngular, omegaDot;
// position = position.rotate(theta, axis);
/*
	The theta used to find the initial position is not the same as the one
	use to rotate the link.
*/
theta = Math.PI/4.0;
// theta = 0.0;
// position = position / numeric.norm2(position);
/*
	theta does not need to be part of the config but it makes my life easiers
	config [x, y, theta, linear velocity, angular velocity]
*/
var configuration = [position[0], position[1], theta, 0, 0, 0, 0, 0, 0];


init();
animate();

function init() 
{

	  scene = new THREE.Scene()
	  renderer = new THREE.CanvasRenderer( )
	  renderer.setSize( width, height)

	camera = new THREE.PerspectiveCamera( 70, window.innerWidth / window.innerHeight, 1, 1000 );
	camera.position.set(1,0,5);
	camera.lookAt( new THREE.Vector3(0,-1,0) );


	var geometry = new THREE.BoxGeometry( b_width, b_height, b_depth );

	var texture = THREE.ImageUtils.loadTexture( 'crate.gif' );
	texture.anisotropy = renderer.getMaxAnisotropy();

	var material = new THREE.MeshBasicMaterial( { map: texture } );
	// var pendulum = new THREE.Object3D();
	
	pendulum = new THREE.Mesh( geometry, material );
	pendulum.position.set(configuration[0], configuration[1], 0);
	pendulum.rotation.set(0,0,configuration[2]);
	scene.add( pendulum );
	// mesh.position.set( 0, -100, 0 );

	ground = new THREE.Mesh( new THREE.PlaneGeometry( 10, 10, 1, 1 ), material );
	ground.position.set( 0, -4, 0 );
	ground.rotation.x = THREE.Math.degToRad(-90);
	scene.add( ground );
	
	centre = new THREE.Mesh( new THREE.SphereGeometry( 0.05, 0.05, 0.05), material );
	centre.position.set( 0, 0, 0);
	scene.add( centre );
	
	/*
	object = new THREE.AxisHelper( 1 );
	object.position.set( 0, 0, 0 );
	scene.add( object );
	*/
	
	drawAcceleration = new THREE.ArrowHelper( new THREE.Vector3( 0, 1, 0 ), new THREE.Vector3( 0, 0, 0 ), 1, 0x0000ff);
	drawAcceleration.position.set( 1, 0, 0 );
	scene.add( drawAcceleration );
	
	drawAngularAcceleration = new THREE.ArrowHelper( new THREE.Vector3( 0, 0, 1 ), new THREE.Vector3( 0, 0, 0 ), 1, 0x00ff00);
	drawAngularAcceleration.position.set( 1, 0, 0 );
	scene.add( drawAngularAcceleration );
	
	
	/*
		Draw velocities of object
		drawVelocity, drawAngularVelocity
	*/
	
	drawVelocity = new THREE.ArrowHelper( new THREE.Vector3( 0, 0, 1 ), new THREE.Vector3( 0, 0, 0 ), 1, 0x00ffff);
	drawVelocity.position.set( 1, 0, 0 );
	scene.add( drawVelocity );

	drawAngularVelocity = new THREE.ArrowHelper( new THREE.Vector3( 0, 0, 1 ), new THREE.Vector3( 0, 0, 0 ), 1, 0xffff00);
	drawAngularVelocity.position.set( 1, 0, 0 );
	scene.add( drawAngularVelocity );
	
	
	
	drawTorqueForce = new THREE.ArrowHelper( new THREE.Vector3( 0, 0, 1 ), new THREE.Vector3( 0, 0, 0 ), 1, 0xff0000);
	drawTorqueForce.position.set( 1, 0, 0 );
	scene.add( drawTorqueForce );


	window.addEventListener( 'resize', onWindowResize, false );
	
	stats = new Stats();
	stats.domElement.style.position = 'absolute';
	stats.domElement.style.top = '0px';
	container.appendChild( stats.domElement );

}

function onWindowResize() {

	camera.aspect = window.innerWidth / window.innerHeight;
	camera.updateProjectionMatrix();

	renderer.setSize( window.innerWidth, window.innerHeight );

}

function animate() {

	requestAnimationFrame( animate );

	// mesh.rotation.x += 0.005;
	// mesh.rotation.y += 0.001;
	
	render(deltaTime, frameNumber, simTime);
	frameNumber = frameNumber+1;
	simTime = simTime + deltaTime;
	stats.update();

}

function animatePendulum( dt )
{
	
}


function render( dt, _frameNumber, _simTime ) {


	// camera.position.x = Math.cos( timer ) * 800;
	// camera.position.z = Math.sin( timer ) * 800;

	// camera.lookAt( scene.position );

	for ( var i = 0, l = scene.children.length; i < l; i ++ ) 
	{

		var object = scene.children[ i ];
		// object.rotation.x = timer * 5;
		// object.rotation.y = timer * 2.5;

	}
	updateState(configuration);
	// printData();
	
	renderer.render( scene, camera );

}

function updateState( config )
{
	pendulum.position.set(config[0], config[1], 0);
	pendulum.rotation.set(0,0,configuration[2]);
}

