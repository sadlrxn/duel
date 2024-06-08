import gsap from "gsap";

export function glowText({
  target,
  delay = 0,
  repeatDelay = 0.3,
}: {
  target: gsap.TweenTarget;
  delay: number;
  repeatDelay: number;
}) {
  return gsap.timeline().fromTo(
    target,
    {
      opacity: 0.3,
    },
    {
      opacity: 1,
      yoyo: true,
      repeat: -1,
      delay,
      duration: 0,
      repeatDelay,
    }
  );
}
