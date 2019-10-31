import Chart from 'chart.js';

export default class Leaderboard {
constructor() {
  const canvas = <HTMLCanvasElement> document.getElementById("graph");
  let chart = new Chart(canvas.getContext("2d"), {
    type: "line",
    options: {
      scales: {
        xAxes: [
          {
            type: "time",
            time: {
              // @ts-ignore
              unit: "minute",
              // @ts-ignore
              max: Date.now(),
              tooltipFormat: "dd HH:mm:ss"
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
              console.log(tooltipItem.datasetIndex);
              // @ts-ignore
              dataset = data.datasets[tooltipItem.datasetIndex];
              // @ts-ignore
              datapoint = dataset.data[tooltipItem.index];
              // @ts-ignore
              challName = datapoint.tooltipLabel;
              // @ts-ignore
              userName = dataset.label;
              // @ts-ignore
              points = datapoint.y;
              // @ts-ignore
              return userName + ': ' + challName + ' ('+points+')';
            }
         }
      },
      mantainAspectRatio: false,
      responsive: true,
      showScale: false
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

    chart.data = {
      datasets: score.chart
    };
    chart.update();

    let table = document.getElementsByTagName("tbody")[0];
    table.innerHTML = "<tr><th>Name</th><th>Points</th></tr>";
    for (let i = 0; i < score.table.name.length; i++) {
      let row = table.insertRow();
      let namecell = row.insertCell();
      namecell.innerHTML = score.table.name[i];
      let pointcell = row.insertCell();
      pointcell.innerHTML = score.table.points[i];
    }
  };
  window.addEventListener("resize", function() {
    chart.update();
  });
  }
}
