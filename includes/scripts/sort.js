function sortTable(n) {
    var table, rows, switching, i, x, y, shouldSwitch, dir, switchcount = 0;
    table = document.getElementById("pensList");
    switching = true;
    // Set the sorting direction to ascending:
    dir = "asc";

    /* Make a loop that will continue until
       no switching has been done:*/
    while (switching) {
        // Start by saying: no switching is done:
        switching = false;
        rows = table.rows;

        /* Loop through all table rows (except the
           first, which contains table headers):*/
        for (i = 1; i < (rows.length - 1); i++) {
            // Start by saying there should be no switching:
            shouldSwitch = false;

            /* Get the two elements you want to compare,
               one from the current row and one from the next:*/
            x = rows[i].getElementsByTagName("TD")[n];
            y = rows[i + 1].getElementsByTagName("TD")[n];

            // Check if the column is a date column (modify as needed):
            if (n === 8) { // Assuming the ninth column is the date column
                var dateX = new Date(x.innerHTML);
                var dateY = new Date(y.innerHTML);

                if ((dir == "asc" && dateX > dateY) || (dir == "desc" && dateX < dateY)) {
                    // If so, mark as a switch and break the loop:
                    shouldSwitch = true;
                    break;
                }
            } else if (n === 0 || n === 10) { // Numeric columns
                if ((dir == "asc" && parseFloat(x.innerHTML) > parseFloat(y.innerHTML)) ||
                    (dir == "desc" && parseFloat(x.innerHTML) < parseFloat(y.innerHTML))) {
                    // If so, mark as a switch and break the loop:
                    shouldSwitch = true;
                    break;
                }
            } else {
                // For non-numeric and non-date columns, perform regular string comparison:
                if (dir == "asc" ? x.innerHTML.toLowerCase() > y.innerHTML.toLowerCase() : x.innerHTML.toLowerCase() < y.innerHTML.toLowerCase()) {
                    shouldSwitch = true;
                    break;
                }
            }
        }

        if (shouldSwitch) {
            // If a switch has been marked, make the switch:
            rows[i].parentNode.insertBefore(rows[i + 1], rows[i]);
            switching = true;
            // Each time a switch is done, increase this count by 1:
            switchcount++;
        } else {
            // If no switching has been done and the direction is "asc":
            if (switchcount == 0 && dir == "asc") {
                // Set the direction to "desc" and run the while loop again:
                dir = "desc";
                switching = true;
            }
        }
    }
}
