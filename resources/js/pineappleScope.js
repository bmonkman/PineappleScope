function deleteFiring(id) {
    var snackbarContainer = document.querySelector('#toast-container');

    var handler = function(e) {
        snackbarContainer.MaterialSnackbar.cleanup_()
        var xhttp = new XMLHttpRequest();
        xhttp.onreadystatechange = function() {
        if (this.readyState == 4 && this.status == 200) {
            var data = {message: 'Deleted firing.'};
            snackbarContainer.MaterialSnackbar.showSnackbar(data);
            setTimeout(function() {document.location = "/";} , 2000)

        } else if (this.readyState == 4) {
            var data = {message: 'Error while trying to delete..'};
            snackbarContainer.MaterialSnackbar.showSnackbar(data);
        }
        };
        xhttp.open("DELETE", "/firing/"+id, true);
        xhttp.send();
    };

    var data = {
        message: 'Are you sure?',
        actionHandler: handler,
        timeout: 5000,
        actionText: 'Yes'
    };
    snackbarContainer.MaterialSnackbar.showSnackbar(data);
}



function renderChart(data){
    var ctx = document.getElementById("myChart").getContext('2d');

    myChart = new Chart(ctx, {
        type: 'line',
        data: data,
        options: {
            tooltips: {
                mode: 'nearest',
                intersect: false
            },
            scales: {
                xAxes: [{
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
                    }],
                yAxes: [{
                    id: 'A',
                    type: 'linear',
                    position: 'left',
                    ticks: {
                        beginAtZero:false
                    }
                }, {
                    id: 'B',
                    type: 'linear',
                    position: 'right',
                    display: false
                }]
            }

        }
    });
}