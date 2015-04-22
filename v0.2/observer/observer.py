from OpenGL.GLUT import *
from OpenGL.GLU import *
from OpenGL.GL import *
import sys
import socket
from threading import Thread
from time import sleep

UDP_IP = "127.0.0.1"
UDP_PORT = 9995
sock = socket.socket(socket.AF_INET, # Internet
                     socket.SOCK_DGRAM) # UDP
sock.bind((UDP_IP, UDP_PORT))

name = 'ARM: Game'

agentDB={}
def main():
    glutInit(sys.argv)
    glutInitDisplayMode(GLUT_DOUBLE | GLUT_RGB | GLUT_DEPTH)
    glutInitWindowSize(800,800)
    glutCreateWindow(name)
    
    thread = Thread(target = agentUpdate, args = (10, ))
    thread.start()
    
    
    glClearColor(0.9,0.9,0.9,1.)
    glShadeModel(GL_SMOOTH)
    glEnable(GL_CULL_FACE)
    glEnable(GL_DEPTH_TEST)
    glEnable(GL_LIGHTING)
    # blending 
    # glEnable(GL_BLEND)
    # glBlendFunc(GL_SRC_ALPHA,GL_ONE)
    
    lightZeroPosition = [10.,4.,10.,1.]
    lightZeroColor = [0.8,1.0,0.8,1.0] #green tinged
    glLightfv(GL_LIGHT0, GL_POSITION, lightZeroPosition)
    glLightfv(GL_LIGHT0, GL_DIFFUSE, lightZeroColor)
    glLightf(GL_LIGHT0, GL_CONSTANT_ATTENUATION, 0.1)
    glLightf(GL_LIGHT0, GL_LINEAR_ATTENUATION, 0.05)
    glEnable(GL_LIGHT0)
    
    glutDisplayFunc(display)
    glMatrixMode(GL_PROJECTION)
    gluPerspective(40.,1.,1.,100.0)
    glMatrixMode(GL_MODELVIEW)
    gluLookAt(20,20,35,
              0,0,0,
              0,1,0)
    glPushMatrix()
    glutMainLoop()
    thread.join()
    return

def display():
    while True:
        glClear(GL_COLOR_BUFFER_BIT|GL_DEPTH_BUFFER_BIT)
        # glClear(GL_COLOR_BUFFER_BIT)
        glPushMatrix()
        color = [0.0,0.0,1.0,0.8]
        # glColor4f(0,0,1,.5)
        glMaterialfv(GL_FRONT,GL_DIFFUSE,color)
        glMaterialfv(GL_FRONT,GL_SPECULAR,color)
        glutWireCube(21)
        glPopMatrix()

        # print "Size of agent DB" + str(len(agentDB))
        for agent, location in agentDB.iteritems():
            glPushMatrix()
            glTranslatef(location[0],location[1], location[2])
            color = [1.0,0.,0.,1.]
            # glColor4f(1,0,0,.5);
            glMaterialfv(GL_FRONT,GL_DIFFUSE,color)
            glMaterialfv(GL_FRONT,GL_SPECULAR,color)
            glutSolidSphere(0.5,20,20)
            glPopMatrix()
    
        glutSwapBuffers()
    return

def agentUpdate(_var):
    import json
    while True:
        data, addr = sock.recvfrom(1024) # buffer size is 1024 bytes
        json_data = json.loads(data)
        # print "received json message:", json_data
        # print "agent:", json_data['Agent']
        # print "Location:", json_data['Location']
        if json_data['Action'] == 'UpdateLocation':
            agentDB[json_data['Agent']] = [float(json_data['Location']['X']),float(json_data['Location']['Y']),float(json_data['Location']['Z'])]

if __name__ == '__main__': main()