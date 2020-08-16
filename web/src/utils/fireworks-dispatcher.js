import fireworks from "fireworks";

const dispatchFireworks = function() {
  fireworks({
    x: window.innerWidth / 2,
    y: window.innerHeight / 2,
    colors: ["#cc3333", "#4CAF50", "#0083F7", "#F2C342"],
    canvasHeight: 700,
    canvasWidth: 700,
    particleTimeout: 3000,
    count: 160,
    bubbleSizeMinimum: 5,
    bubbleSizeMaximum: 20,
    bubbleSpeedMinimum: 8,
    bubbleSpeedMaximum: 15,
  });
};

export default dispatchFireworks;
