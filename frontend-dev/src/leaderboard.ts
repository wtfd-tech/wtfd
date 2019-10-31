import Chart from 'chart.js';
import 'moment';

export default class Leaderboard {
constructor() {
  const canvas = <HTMLCanvasElement> document.getElementById("graph");
  let chart = new Chart(canvas.getContext("2d"), {
    type: "line",
    options: {
      spanGaps: true,
      scales: {
        xAxes: [
          {
            type: "time",
            time: {
              parser: "ddd MMM DD HH:mm:ss ZZ Y",
              // @ts-ignore
              unit: "minute",
              tooltipFormat: "ddd HH:mm:ss"
            },
            ticks: {
              source: "auto"
            }
          }
        ],
        yAxes: [
          {
            ticks: {
              beginAtZero: true
            }
          }
        ]
      },
      elements: {
        line: {
          tension: 0,
          fill: false
        }
      },
      legend: {
        display: false
      },
      tooltips: {
         callbacks: {
            label: function(tooltipItem, data) {
              let dataset = data.datasets[tooltipItem.datasetIndex];
              let datapoint = dataset.data[tooltipItem.index];
              // @ts-ignore
              let challName = datapoint.tooltipLabel;
              let userName = dataset.label;
              // @ts-ignore
              let points = datapoint.y;
              return userName + ': ' + challName + ' ('+points+')';
            }
         }
      },
      maintainAspectRatio: false,
      responsive: true,
//      showScale: false
    }
  });

  let ws = new WebSocket("ws://" + window.location.host + "/ws");
  ws.onopen = function() {
    // Web Socket is connected, send data uting send()
    console.log("ws connected");
  };
  ws.onclose = function() {
    alert("WS Disconnected, reload the page");
  };

  ws.onmessage = evt => {
    var rec = evt.data;
    console.log(evt.data);
    let score = JSON.parse(rec);

    let table = <HTMLTableElement> document.getElementById("leaderboardtable");
    table.innerHTML = "<tr><th>Name</th><th>Points</th></tr>";
    for (let i = 0; i < score.table.name.length; i++) {
      let row = table.insertRow();
      let namecell = row.insertCell();
      namecell.innerHTML = score.table.name[i];
      let pointcell = row.insertCell();
      pointcell.innerHTML = score.table.points[i];
    }

    chart.data = {
      datasets: score.chart,
      labels: ["a"]
    };
    // @ts-ignore
    chart.options.scales.xAxes[0].ticks.min = chart.data.datasets[0].data[0].t;
    chart.options.scales.xAxes[0].ticks.max = Date.now();
    chart.update();

  };
  window.addEventListener("resize", function() {
    chart.update();
  });
  }
}
