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