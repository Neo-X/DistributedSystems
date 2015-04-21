fs = require("fs")

THREE = require("three.js")
join = require("path").join

// app.get '/test/top_:top_id/side_:side_id/x_:x/y_:y.jpg', (req, res, next) =>

  var width = 660
  var height = 500

  camera = new THREE.PerspectiveCamera(50, width / height, 1, 1000)
  scene = new THREE.Scene()
  renderer = new THREE.CanvasRenderer( )
  renderer.setSize( width, height)

  camera.position.z = 100

  camera_container = new THREE.Object3D
  scene.add( camera_container)
  camera_container.add( camera)

  camera.position.z = 75

  // We have one background plane
  plane_image = new Image()
  // plane_image.src = fs.readFileSync TOP_DIR + "public/images/vtx_logo.jpg"
  texture = new THREE.Texture( plane_image, new THREE.UVMapping())
  texture.needsUpdate = true

  loader = new THREE.JSONLoader()

  geometry = new THREE.PlaneGeometry(200, 200)
  material = new THREE.MeshBasicMaterial
    color       : 0x698144
    // #shading        : THREE.SmoothShading
    map     : texture
    overdraw: true
  plane = new THREE.Mesh( geometry, material)
  plane.position.z = -50
  plane.position.y = -4
  plane.position.x = 4.5

  // # We also have an object in the foreground
  scene.add( plane)
  geometry = false
  loader.createModel( JSON.parse(fs.readFileSync(TOP_DIR + 'public/blender_export.json'))), (done) =>
    geometry = done

  # Imager.texture gives us a canvas based on some code that grabs specific info
  texture = new THREE.Texture (Imager.texture req.params.side_id, req.params.top_id), new THREE.UVMapping()
  texture.needsUpdate = true

  material = new THREE.MeshBasicMaterial
    color: 0xaaaaaa
    map: texture
    overdraw: true
  mesh = new THREE.Mesh geometry, material

  mesh.rotation.x = parseFloat req.params.x
  mesh.rotation.y = parseFloat req.params.y

  scene.add mesh
  mesh.dynamic = true
  renderer.render scene, camera

  renderer.domElement.toBuffer (err, buf) ->
    res.contentType 'image/jpg'
    res.send buf