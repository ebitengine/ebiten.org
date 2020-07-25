package main

var Time float
var Cursor vec2

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
	lightpos := vec3(Cursor, 50)
	lightdir := normalize(lightpos - position.xyz)
	normal := normalize(texture1At(texCoord) - 0.5)
	ambient := 0.25
	diffuse := 0.75 * max(0.0, dot(normal.xyz, lightdir))
	return texture0At(texCoord) * (ambient + diffuse)
}
