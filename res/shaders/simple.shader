#shader vertex
#version 330

uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;

in vec3 vert;
in vec2 vertTexCoord;

out vec2 fragTexCoord;

void main() {
    fragTexCoord = vertTexCoord;
    gl_Position = projection * camera * model * vec4(vert, 1); 
}

#shader fragment
#version 330

in vec2 fragTexCoord;

out vec4 outputColor;

uniform sampler2D tex;

void main() {
    outputColor = texture(tex, fragTexCoord);
}