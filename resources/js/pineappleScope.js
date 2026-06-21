// Tint the header from cold (dark blue) to hot (orange) based on kiln temp.
// Interpolates in HSL so the hue sweeps the warm way: blue -> purple -> red -> orange.
function tempToColor(tempC) {
    var COLD = 20, HOT = 1000;
    var t = Math.max(0, Math.min(1, (tempC - COLD) / (HOT - COLD)));
    var h = (210 + t * 176) % 360;   // 210 (blue) -> 386 wraps to 26 (orange)
    var s = 29 + t * 70;             // 29% -> 99%
    var l = 24 + t * 36;             // 24% -> 60%
    return 'hsl(' + h.toFixed(0) + ' ' + s.toFixed(0) + '% ' + l.toFixed(0) + '%)';
}

function applyHeaderTemp() {
    var header = document.querySelector('.app-header');
    if (!header) return;
    var temp = parseFloat(header.getAttribute('data-temp'));
    if (isNaN(temp)) return;
    header.style.background = tempToColor(temp);
}

var toastHideTimer;

function showToast(message, actionText, actionHandler, timeout) {
    var container = document.getElementById('toast-container');
    var msg = document.getElementById('toast-message');
    var action = document.getElementById('toast-action');

    // Cancel any pending hide so a previous toast's timer can't dismiss this one.
    clearTimeout(toastHideTimer);

    msg.textContent = message;

    if (actionText) {
        action.textContent = actionText;
        action.style.display = '';
        action.onclick = function () {
            container.classList.remove('toast--show');
            actionHandler();
        };
    } else {
        action.style.display = 'none';
        action.onclick = null;
    }

    container.classList.add('toast--show');

    if (timeout) {
        toastHideTimer = setTimeout(function () { container.classList.remove('toast--show'); }, timeout);
    }
}

function deleteFiring(id) {
    var doDelete = function () {
        var xhttp = new XMLHttpRequest();
        xhttp.onreadystatechange = function () {
            if (this.readyState == 4 && this.status == 200) {
                showToast('Deleted firing.', null, null, 2000);
                setTimeout(function () { document.location = "/"; }, 1500);
            } else if (this.readyState == 4) {
                showToast('Error while trying to delete..', null, null, 4000);
            }
        };
        xhttp.open("DELETE", "/firing/" + id, true);
        xhttp.send();
    };

    showToast('Are you sure?', 'Yes', doDelete, 5000);
}



function renderChart(data){
    var ctx = document.getElementById("myChart").getContext('2d');

    myChart = new Chart(ctx, {
        type: 'line',
        data: data,
        options: {
            responsive: true,
            maintainAspectRatio: false,
            tooltips: {
                mode: 'nearest',
                intersect: false
            },
            scales: {
                x: {
                    type: "time",
                    time: {
                        unit: 'minute',
                        unitStepSize: 30,
                        round: 'minute',
                        tooltipFormat: "h:mm:ss a",
                        displayFormats: {
                        hour: 'MMM D, h:mm'
                        }
                    }
                },
                inner: {
                    type: 'linear',
                    position: 'left',
                    beginAtZero:false
                },
                outer: {
                    type: 'linear',
                    position: 'right',
                    display: false,
                    beginAtZero:false
                }
            },
            plugins: {
                zoom: {
                  zoom: {
                    wheel: {
                      enabled: true,
                    },
                    pinch: {
                      enabled: true
                    },
                    mode: 'x',
                  },
                  pan: {
                    enabled: true,
                    mode: 'x',
                  },
                  limits: {
                    x: {max: 'original', min: 'original'},
                  },
                }
            }
        }
    });
}