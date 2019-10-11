var chart;
(function() {
  var c = document.getElementById("graph").getContext("2d");

  chart = new Chart(c, {
    type: "line",
    options: {
      scales: {
        xAxes: [
          {
            type: "time",
            time: {
              unit: "minute",
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
              dataset = data.datasets[tooltipItem.datasetIndex];
              datapoint = dataset.data[tooltipItem.index];
              challName = datapoint.tooltipLabel;
              userName = dataset.label;
              points = datapoint.y;
              return userName + ': ' + challName + ' ('+points+')';
            }
         }
      },
      mantainAspectRatio: false,
      responsive: true,
      showScale: false
    }
  });

  ws = new WebSocket("ws://" + window.location.host + "/ws");
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
    score = JSON.parse(rec);

    chart.data = {
      datasets: score.chart
    };
    chart.update();

    table = document.getElementsByTagName("tbody")[0];
    table.innerHTML = "<tr><th>Name</th><th>Points</th></tr>";
    for (i = 0; i < score.table.name.length; i++) {
      row = table.insertRow();
      namecell = row.insertCell();
      namecell.innerHTML = score.table.name[i];
      pointcell = row.insertCell();
      pointcell.innerHTML = score.table.points[i];
    }
  };
  window.addEventListener("resize", function() {
    chart.update();
  });
})();
