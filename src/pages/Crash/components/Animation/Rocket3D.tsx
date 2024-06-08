import { useRef } from 'react';
import styled from 'styled-components';
import useSpline from '@splinetool/r3f-spline';
import { Canvas, useFrame } from '@react-three/fiber';
import { PerspectiveCamera } from '@react-three/drei';
import { Group } from 'three';

const splinePath = window.origin + '/assets/scene.splinecode';

function Scene({ speed = 1 }: any) {
  const { nodes, materials } = useSpline(splinePath);

  const rocketRef = useRef<Group>(null);

  useFrame((_state, delta) => {
    if (!rocketRef || !rocketRef.current) return;
    rocketRef.current.rotation.y += speed * delta * 4;
  });

  return (
    <>
      <group dispose={null}>
        <PerspectiveCamera
          name="Camera"
          makeDefault={true}
          zoom={0.5}
          far={100000}
          near={70}
          fov={45}
          position={[0, 0, 1000]}
          rotation={[0, 0, 0]}
        />
        <group name="Deleted Components" visible={false}>
          <group name="Fins" position={[3.96, 205.58, 0]} scale={[1.1, 1, 1]}>
            <mesh
              name="Fins1"
              geometry={nodes.Fins1.geometry}
              material={materials['Fins1 Material']}
              castShadow
              receiveShadow
            />
          </group>
        </group>
        <group ref={rocketRef} name="Position" position={[0, 0, 0]}>
          <group name="Yaw" position={[0, 0, 0]} rotation={[0, 0, 0]} scale={1}>
            <group name="Rocket" position={[0, -0.01, 0]}>
              <group name="Fire Light" position={[0, -156.62, -3.26]}>
                <pointLight
                  name="Fire Light1"
                  intensity={2.5}
                  distance={500}
                  shadow-mapSize-width={1024}
                  shadow-mapSize-height={1024}
                  shadow-camera-near={100}
                  shadow-camera-far={100000}
                  color="#befbfe"
                  position={[0, -69.25, 134.06]}
                />
                <pointLight
                  name="Fire Light 2"
                  intensity={2.5}
                  distance={500}
                  shadow-mapSize-width={1024}
                  shadow-mapSize-height={1024}
                  shadow-camera-near={100}
                  shadow-camera-far={100000}
                  color="#befbfe"
                  position={[0, -69.25, -115.39]}
                />
              </group>
              <group name="Fire Trail" position={[0, -259.41, 0]}>
                <mesh
                  name="Blurry Fire Trail"
                  geometry={nodes['Blurry Fire Trail'].geometry}
                  material={materials['Blurry Fire Trail Material']}
                  castShadow
                  receiveShadow
                  position={[0, 22.23, 0]}
                  scale={[2.7, 3.85, 2.7]}
                />
                <mesh
                  name="Fire Trail1"
                  geometry={nodes['Fire Trail1'].geometry}
                  material={materials['Fire Trail1 Material']}
                  castShadow
                  receiveShadow
                  position={[0, 47.3, 0]}
                  scale={[2, 3, 2]}
                />
              </group>
              <mesh
                name="Window"
                geometry={nodes.Window.geometry}
                material={materials['Window Material']}
                castShadow
                receiveShadow
                position={[-0.9, 267.15, 65.62]}
                rotation={[0, -0.08, 0]}
                scale={[1.61, 1.14, 1.14]}
              />
              <group name="Engine" position={[0, -46.37, 0]}>
                <mesh
                  name="4"
                  geometry={nodes['4'].geometry}
                  material={materials['4 Material']}
                  castShadow
                  receiveShadow
                  position={[0, -51.58, 0]}
                  rotation={[Math.PI, 0, 0]}
                />
                <mesh
                  name="3"
                  geometry={nodes['3'].geometry}
                  material={materials['3 Material']}
                  castShadow
                  receiveShadow
                  position={[0, -26.84, 0]}
                  rotation={[Math.PI, 0, 0]}
                />
                <mesh
                  name="2"
                  geometry={nodes['2'].geometry}
                  material={materials['2 Material']}
                  castShadow
                  receiveShadow
                  position={[0, -4.05, 0]}
                  rotation={[Math.PI, 0, 0]}
                />
                <mesh
                  name="1"
                  geometry={nodes['1'].geometry}
                  material={materials['1 Material']}
                  castShadow
                  receiveShadow
                  position={[0, 20.68, 0]}
                  rotation={[Math.PI, 0, 0]}
                />
              </group>
              <group
                name="Fins Clones"
                position={[3.96, 205.58, 0]}
                scale={[1.1, 1, 1]}
              >
                <group name="Clone 0" position={[50, 0, 0]}>
                  <mesh
                    name="Fins2"
                    geometry={nodes.Fins2.geometry}
                    material={materials['Untitled Material']}
                    castShadow
                    receiveShadow
                  />
                </group>
                <group
                  name="Clone 1"
                  position={[0, 0, 50]}
                  rotation={[0, -Math.PI / 2, 0]}
                >
                  <mesh
                    name="Fins3"
                    geometry={nodes.Fins3.geometry}
                    material={materials['Untitled Material']}
                    castShadow
                    receiveShadow
                  />
                </group>
                <group
                  name="Clone 2"
                  position={[-50, 0, 0]}
                  rotation={[-Math.PI, 0, -Math.PI]}
                >
                  <mesh
                    name="Fins4"
                    geometry={nodes.Fins4.geometry}
                    material={materials['Untitled Material']}
                    castShadow
                    receiveShadow
                  />
                </group>
                <group
                  name="Clone 3"
                  position={[0, 0, -50]}
                  rotation={[0, Math.PI / 2, 0]}
                >
                  <mesh
                    name="Fins5"
                    geometry={nodes.Fins5.geometry}
                    material={materials['Untitled Material']}
                    castShadow
                    receiveShadow
                  />
                </group>
              </group>
              <mesh
                name="Body"
                geometry={nodes.Body.geometry}
                material={materials['Body Material']}
                castShadow
                receiveShadow
                position={[0, 198.32, 0]}
                scale={5}
              />
            </group>
          </group>
        </group>
        <directionalLight
          name="Directional Light"
          castShadow
          intensity={2}
          shadow-mapSize-width={2048}
          shadow-mapSize-height={2048}
          shadow-camera-near={-10000}
          shadow-camera-far={100000}
          shadow-camera-left={-1000}
          shadow-camera-right={1000}
          shadow-camera-top={1000}
          shadow-camera-bottom={-1000}
          position={[-454.76, 663.87, 0]}
        />
        <hemisphereLight
          name="Default Ambient Light"
          intensity={0.75}
          color="#eaeaea"
        />
      </group>
    </>
  );
}

export default function Rocket3D({ speed = 1, ...props }: any) {
  return (
    <CustomCanvas flat linear {...props}>
      <Scene speed={speed} />
    </CustomCanvas>
  );
}

const CustomCanvas = styled(Canvas)`
  canvas {
    position: absolute;
    left: 50%;
    top: 50%;
    transform: translate(-50%, -50%);
  }
`;
