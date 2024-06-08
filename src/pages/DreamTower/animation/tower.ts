import gsap from "gsap";
import { glowText } from "./text";
import { rowFlowAnims, setRowBackground } from "./row";
import { textSpanIntersectsWithTextSpan } from "typescript";

export function towerResetAnimation({
  textTargets,
  rowTargets1,
  rowTargets2,
  rowTargets3,
  rowTargets4,
}: {
  textTargets: gsap.TweenTarget[];
  rowTargets1: gsap.TweenTarget[];
  rowTargets2: gsap.TweenTarget[];
  rowTargets3: gsap.TweenTarget[];
  rowTargets4: gsap.TweenTarget[];
}) {
  const resetTimeline = gsap.timeline();
  textTargets.forEach((target) => {
    target && resetTimeline.add(gsap.set(target, { opacity: 1 }));
  });
  rowTargets1.forEach((target) => {
    target &&
      resetTimeline.add(
        gsap.set(target, {
          background: "#73218F00",
          borderLeftColor: "transparent",
          borderRightColor: "transparent",
        })
      );
  });
  rowTargets2.forEach((target) => {
    target &&
      resetTimeline.add(
        gsap.set(target, {
          background:
            "linear-gradient(90deg, #73218Fa0 0%, #1A293D00 50%, #73218Fa0 100%)",
          borderLeftColor: "#e8caee",
          borderRightColor: "#e8caee",
        })
      );
  });
  rowTargets3.forEach((target) => {
    target &&
      resetTimeline.add(
        gsap.set(target, {
          background:
            "linear-gradient(90deg, #8F2121a0 0%, #1A293D00 50%, #8F2121a0 100%)",
          borderLeftColor: "#e8caee",
          borderRightColor: "#e8caee",
        })
      );
  });
  rowTargets4.forEach((target) => {
    target &&
      resetTimeline.add(
        gsap.set(target, {
          background:
            "linear-gradient(90deg, #2F23B7a0 0%, #1A293D00 50%, #2F23B7a0 100%)",
          borderLeftColor: "#e8caee",
          borderRightColor: "#e8caee",
        })
      );
  });
  return resetTimeline;
}

export function towerWinAnimation({
  textTargets,
  rowTargets,
  duration = 0.2,
}: {
  textTargets: gsap.TweenTarget[];
  rowTargets: gsap.TweenTarget[];
  duration?: number;
}) {
  const rowFlowTl = gsap.timeline();
  rowFlowTl.add(rowFlowAnims({ targets: rowTargets, duration }));

  const glowTextTl = gsap.timeline({ pause: true });
  textTargets.forEach((target, index) => {
    target &&
      glowTextTl.add(
        glowText({ target, delay: duration * index, repeatDelay: duration }),
        0
      );
  });

  const opacityValue = {
    value: 0.51,
  };

  const tl = gsap.timeline({
    onStart: () => {
      glowTextTl.progress(0).play();
    },
  });

  tl.add(glowTextTl, 0)
    .add(rowFlowTl, 0)
    .add(
      gsap.delayedCall(rowFlowTl.totalDuration(), () => {
        glowTextTl.pause();
      }),
      0
    )
    .fromTo(
      opacityValue,
      {
        value: 0.51,
      },
      {
        value: 2.49,
        roundProps: "value",
        repeat: 1,
        duration: 1,
        repeatDelay: 0.5,
        onUpdate: () => {
          const dark = opacityValue.value === 1;
          textTargets.forEach((target) => {
            target && gsap.set(target, { opacity: dark ? 0.3 : 1 });
          });
          rowTargets.forEach((target) => {
            target && setRowBackground({ target, left: !dark, right: !dark });
          });
        },
      },
      ">"
    )
    .add(
      gsap.delayedCall(3, () => {
        tl.progress(0);
      })
    );

  return tl;
}
