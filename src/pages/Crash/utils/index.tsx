const MAX_SPEED = 3;

export const calculateSpeed = (timeElapsed: number) => {
  // let speed = Math.pow(1.0373, timeElapsed);
  let speed = Math.pow(1.02, timeElapsed);
  // let speed = 1 + (timeElapsed / 100) * 2;
  if (speed > MAX_SPEED) speed = MAX_SPEED;
  return speed;
};

export const calculateAngle = (timeElapsed: number) => {
  const speed = calculateSpeed(timeElapsed);
  const angle = 45 + (30 * (speed - 1)) / (MAX_SPEED - 1);
  return angle;
};

export const calculateAngleWithSpeed = (speed: number) => {
  const angle = 45 + (30 * (speed - 1)) / (MAX_SPEED - 1);
  return angle;
};
