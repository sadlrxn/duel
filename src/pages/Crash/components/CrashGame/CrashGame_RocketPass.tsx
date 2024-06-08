import React, {
  useEffect,
  useCallback,
  useState,
  useRef,
  useMemo
} from 'react';
import styled from 'styled-components';
import gsap from 'gsap';
import { Chart } from 'react-chartjs-2';
import 'chart.js/auto';

import { Flex, Span } from 'components';
import { useAppSelector } from 'state';
import { getChartGradient } from 'utils/getThemeValue';

import { Rocket } from '../Animation';

const calculateMultiplier = (timeElapsed: number) => {
  return 1.0024 * Math.pow(1.0718, timeElapsed);
};

interface CrashGameProps {
  setDirection?: any;
}

export default function CrashGame({ setDirection }: CrashGameProps) {
  const { status, time: globalTime } = useAppSelector(state => state.crash);

  const [chartData, setChartData] = useState<any>({ datasets: [] });
  const [chartOptions, setChartOptions] = useState({});
  const [liveMultiplier, setLiveMultiplier] = useState('CONNECTING...');
  const [timeMax, setTimeMax] = useState(0);
  const [chartSwitch, setChartSwitch] = useState(false);
  const [leftTime, setLeftTime] = useState(0);
  const [position, setPosition] = useState({ x: 0, y: 0 });
  const [angle, setAngle] = useState(0);
  const [size, setSize] = useState({ width: 0, height: 0 });

  const [multiMax, setMultiMax] = useState(2.5);

  const timeCount_xaxis = useRef<number[]>([]);
  const multiplierCount = useRef<number[]>([]);
  const graphRef = useRef<HTMLDivElement>(null);

  const resizeObserver = useMemo(
    () =>
      new ResizeObserver(() => {
        if (!graphRef || !graphRef.current) return;
        const width = graphRef.current.clientWidth;
        const height = graphRef.current.clientHeight;
        setSize({ width, height });
      }),
    [graphRef]
  );

  useEffect(() => {
    if (!graphRef || !graphRef.current) return;
    resizeObserver.observe(graphRef.current);
  }, [resizeObserver]);

  useEffect(() => {
    setChartSwitch(true);
    if (!graphRef || !graphRef.current) return;
    const width = graphRef.current.clientWidth;
    const height = graphRef.current.clientHeight;
    setSize({ width, height });
  }, []);

  const calculate = useCallback(() => {
    let timeElapsed = (Date.now() - globalTime) / 1000.0;
    let liveMultiplier = calculateMultiplier(timeElapsed);

    setLiveMultiplier(liveMultiplier.toFixed(2));

    if (multiplierCount.current.length < 1) {
      multiplierCount.current = multiplierCount.current.concat([1]);
      timeCount_xaxis.current = timeCount_xaxis.current.concat([0]);
    }

    let timeMax = Math.max((timeElapsed / 5) * 6, 18);
    setTimeMax(timeMax);
    multiplierCount.current = multiplierCount.current.concat([liveMultiplier]);
    timeCount_xaxis.current = timeCount_xaxis.current.concat([timeElapsed]);

    let multiMax = calculateMultiplier((18 / 6) * 5);
    multiMax = Math.max(liveMultiplier, multiMax);
    multiMax = multiMax * 1.1;
    setMultiMax(multiMax);

    if (multiplierCount.current.length > 2) {
      const xMin = 0,
        xMax = timeMax;
      const yMin = 1,
        yMax = multiMax;
      const length = multiplierCount.current.length;
      const width = size.width - 50;
      const height = size.height - 20;

      let pos1 = {
          x: timeCount_xaxis.current[length - 2],
          y: multiplierCount.current[length - 2]
        },
        pos2 = {
          x: timeCount_xaxis.current[length - 1],
          y: multiplierCount.current[length - 1]
        };

      pos1.x = ((pos1.x - xMin) / (xMax - xMin)) * width;
      pos1.y = ((pos1.y - yMin) / (yMax - yMin)) * height;
      pos1.y = height - pos1.y + 10;

      pos2.x = ((pos2.x - xMin) / (xMax - xMin)) * width;
      pos2.y = ((pos2.y - yMin) / (yMax - yMin)) * height;
      pos2.y = height - pos2.y + 10;

      let speed = { x: pos2.x - pos1.x, y: pos2.y - pos1.y };
      const vec = Math.sqrt(speed.x ** 2 + speed.y ** 2);
      speed.x = (speed.x / vec) * 10;
      speed.y = (speed.y / vec) * 10;

      setAngle(Math.atan2(speed.y, speed.x) + Math.PI / 2);

      setDirection(speed);
      setPosition(pos2);
    }
  }, [globalTime, size, setDirection]);

  useEffect(() => {
    let interval: NodeJS.Timer | null = null;
    let tl: gsap.core.Timeline | null = null;

    switch (status) {
      case 'bet':
        multiplierCount.current = [];
        timeCount_xaxis.current = [];
        setLiveMultiplier('1.00');

        setAngle(Math.PI / 2);
        setDirection({ x: 0, y: 0 });
        setPosition({ x: 0, y: size.height });

        let leftTime = 5 - Math.floor((Date.now() - globalTime) / 1000);
        if (leftTime < 0) leftTime = 0;
        setLeftTime(leftTime);

        interval = setInterval(() => {
          leftTime = 5 - Math.floor((Date.now() - globalTime) / 1000);
          if (leftTime < 0) leftTime = 0;
          setLeftTime(leftTime);
        }, 200);
        break;
      case 'ready':
        if (!graphRef || !graphRef.current) break;
        let dx = 1;
        let dy = 1 - calculateMultiplier(-1);
        dx = (dx / 18) * (size.width - 50);
        dy =
          (dy / (calculateMultiplier((18 / 6) * 5) - 1)) * (size.height - 20);
        let scale = Math.ceil(graphRef.current.offsetLeft / dx);

        let roundValue = {
          value: 0
        };

        let progress = (Date.now() - globalTime) / 1000;
        if (progress > 1) progress = 0.998;
        if (progress < 0) progress = 0;

        setAngle(Math.atan2(dy, -dx) - Math.PI / 2);

        tl = gsap
          .timeline()
          .fromTo(
            roundValue,
            { value: 0 },
            {
              value: 10 ** 3,
              roundProps: 'value',
              duration: 1,
              onUpdate: () => {
                const percent = (1 - roundValue.value / 10 ** 3) * scale;
                setPosition({
                  x: 0 - dx * percent,
                  y: size.height - 10 + dy * percent
                });
              }
            }
          )
          .progress(progress);
        break;
      case 'play':
        setLiveMultiplier('1.00');

        calculate();
        interval = setInterval(() => {
          calculate();
        }, 50);
        break;
      case 'explosion':
        break;
      case 'back':
        break;
    }

    return () => {
      interval && clearInterval(interval);
      tl && tl.clear();
    };
  }, [calculate, globalTime, setDirection, size, status]);

  const sendToChart = useCallback(() => {
    setChartData({
      labels: timeCount_xaxis.current,
      datasets: [
        {
          data: multiplierCount.current,
          fill: false,
          borderWidth: 6.61,
          borderColor: function (context: any) {
            const chart = context.chart;
            const { ctx, chartArea } = chart;

            if (!chartArea) return;
            return getChartGradient(
              ctx,
              chartArea,
              {
                left: 0,
                right: position.x,
                top: position.y,
                bottom: size.height
              },
              [
                {
                  percent: 2.91,
                  color: 'rgba(180, 255, 255, 0.25)'
                },
                {
                  percent: 86.55,
                  color: 'rgba(255, 255, 255, 0)'
                }
              ]
            );
          },
          segment: {
            borderColor: function (context: any) {
              if (context.p1DataIndex >= multiplierCount.current.length - 5)
                return 'transparent';
              return undefined;
            }
          },
          color: 'rgba(255, 255, 255,1)',

          pointRadius: 0,
          lineTension: 0.4
        },
        {
          data: multiplierCount.current,
          fill: false,
          borderWidth: 14.86,
          borderColor: function (context: any) {
            const chart = context.chart;
            const { ctx, chartArea } = chart;

            if (!chartArea) return;
            return getChartGradient(
              ctx,
              chartArea,
              {
                left: 0,
                right: position.x,
                top: position.y,
                bottom: size.height
              },
              [
                {
                  percent: 2.91,
                  color: 'rgba(0, 194, 255, 0.1)'
                },
                {
                  percent: 86.55,
                  color: 'rgba(255, 255, 255, 0)'
                }
              ]
            );
          },
          segment: {
            borderColor: function (context: any) {
              if (context.p1DataIndex >= multiplierCount.current.length - 5)
                return 'transparent';
              return undefined;
            }
          },
          color: 'rgba(255, 255, 255,1)',

          pointRadius: 0,
          lineTension: 0.4
        }
      ]
    });

    setChartOptions({
      events: [],
      maintainAspectRatio: false,
      legend: {
        display: false
      },
      layout: {
        padding: 0
      },
      elements: {
        line: {
          tension: 0.4,
          borderJoinStyle: 'round'
        }
      },
      scales: {
        y: {
          afterFit: function (axis: any) {
            axis.width = 50;
          },
          type: 'linear',
          position: 'right',

          title: {
            display: false,
            text: 'value'
          },
          min: 1,
          max: multiMax,
          ticks: {
            stepSize: 0.025,
            padding: 5,

            color: 'rgba(255, 255, 255,1)',
            maxTicksLimit: 46,
            font: {
              family: 'Inter',
              size: 10,
              weight: 700
            },
            callback: function (value: string, index: number) {
              if (index % 4 !== 0) return '';
              if ((+value * 1000) % 25 !== 0) return '';
              return (+value).toFixed(1);
            }
          },
          border: {
            width: 0
          },
          grid: {
            drawBorder: false,
            drawTicks: true,
            drawOnChartArea: false,
            color: function (context: any) {
              if ((context.tick.value * 1000) % 25 !== 0) return '#fff0';
              if (context.index % 4 === 0) return '#fff';
              return '#fff6';
            }
          }
        },
        x: {
          display: false,
          type: 'linear',
          min: 0,
          max: timeMax
        }
      },
      plugins: {
        legend: { display: false }
      },
      animation: false
    });
  }, [multiMax, timeMax, position, size.height]);

  useEffect(() => {
    const temp_interval = setInterval(() => {
      setChartSwitch(false);
      sendToChart();
    }, 10);

    return () => {
      clearInterval(temp_interval);
      setChartSwitch(true);
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [chartSwitch]);

  return (
    <Container
      ref={graphRef}
      justifyContent={status !== 'bet' ? 'end' : 'center'}
      alignItems={status !== 'bet' ? 'start' : 'center'}
    >
      {status !== 'bet' && (
        <>
          <Chart type="line" data={chartData} options={chartOptions} />
          <Multiplier top={position.y}>
            {liveMultiplier + 'x'}
            <Flex
              width="100%"
              height="1px"
              background="linear-gradient(90deg, rgba(77, 169, 255, 0) 0%, rgba(77, 169, 255, 0.5) 11.98%, rgba(77, 169, 255, 0) 100%)"
            />
          </Multiplier>
        </>
      )}
      {status === 'bet' && (
        <Span color="white" fontWeight={700} fontSize="70px">
          {leftTime}
        </Span>
      )}
      <Rocket
        className="crash_rocket"
        position="absolute"
        left={position.x}
        top={position.y}
        width={100}
        height={160}
        angle={angle}
        explosion={status === 'explosion'}
        visible={status !== 'bet' && status !== 'back'}
      />
    </Container>
  );
}

const Multiplier = styled(Flex)`
  position: absolute;
  right: 0;
  width: 100%;
  align-items: center;
  justify-content: end;

  height: 1px;
  max-height: 1px;

  font-weight: 700;
  font-size: 42px;
  color: white;
`;

const Container = styled(Flex)`
  /* position: absolute;
  right: 0;
  top: 0;
  width: 70%;
  height: 90%;
  z-index: -1; */
  position: relative;
  min-height: 500px;
  height: 500px;
  z-index: -1;
`;
