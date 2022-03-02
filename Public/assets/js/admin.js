// Scheduler
function runAtTime(method, hour, minute, second) {
    (function loop() {
        var now = new Date();
        if (now.getHours() === hour && now.getMinutes() === minute && now.getSeconds() === second) {
            method();
        }
        now = new Date();
        var delay = 60000 - (now % 60000);
        setTimeout(loop, delay);
    })();
}

// CSV Table fetcher
function getTable() {
    $.ajax({
        type: "GET",
        url: "/GetAdminCSV/",
        success: function (data) {
            $('table').replaceWith(arrayToTable(Papa.parse(data).data))
        }
    })
}

function SearchStudent() {
    var name = document.getElementById('name').value;
    document.getElementById('name').value = "";
    $.ajax({
        type: "POST",
        url: "/search/" + btoa(name),
        success: function (data) {
            $('table').replaceWith(arrayToTable(Papa.parse(data).data))
        }
    })
}

// Run this automatically on page load
getTable()

//Run at midnight
runAtTime(getTable, 0, 0, 0)

function arrayToTable(tableData) {
    var table = $('<table class="table"></table>');
    $(tableData).each(function (i, rowData) {
        var row = $('<tr></tr>');
        $(rowData).each(function (j, cellData) {
            if(cellData.length >= 1) {
                row.append($('<td>'+cellData+'</td>'))
            }
        })
        table.append(row)
    })
    return table
}