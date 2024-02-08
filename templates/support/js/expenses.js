document.addEventListener('DOMContentLoaded', function () {

    Get();
        document.getElementById("PostForm1").addEventListener("submit", function (event) {
            event.preventDefault();
            Add(1);
            document.getElementById("PostForm1").reset(); // clear form after submission
            plusModal.style.display = "none"; // close modal after submission
        });

        document.getElementById("PostForm2").addEventListener("submit", function (event) {
            event.preventDefault();
            Add(-1);
            document.getElementById("PostForm2").reset(); // clear form after submission
            minModal.style.display = "none"; // close modal after submission
        });

        // ajax request func
        function Get(callback) {
            var xhr = new XMLHttpRequest();
            xhr.open('GET', 'moneybuddy-production.up.railway.app/getuserdata', true);
            xhr.onload = function () {
                if (this.status == 200) {
                    var u = JSON.parse(this.responseText);
                    var output = '';
                    output += '<tr> <td>' + '<b>Date</b>' + '</td> <td>' + '<b>Amount</b>' + '</td> <td>' + '<b>Category</b>' + ' </td> </tr><br>';
                    for (let i = 0; i < u.Transactions.length; i++) {
                        output += '<tr> <td>' + formatTransactionTime(u.Transactions[i].TransactionTime) + '</td> <td>' + u.Transactions[i].Amount + '</td> <td>' + u.Transactions[i].Category + ' </td> </tr><br>';

                        // function that changes output time
                        function formatTransactionTime(dateTimeString) {
                            var dateTime = new Date(dateTimeString);
                            var month = ('0' + (dateTime.getMonth() + 1)).slice(-2);
                            var day = ('0' + dateTime.getDate()).slice(-2);
                            var hours = ('0' + dateTime.getHours()).slice(-2);
                            var minutes = ('0' + dateTime.getMinutes()).slice(-2);
                            return day + '-' + month + ' ' + hours + ':' + minutes;
                        }
                        
                    }
                    document.getElementById("transact").innerHTML = output;

                }
            };
            xhr.send();
        }
        
        // add values to database
        function Add(value) {
                if (value == 1) {
                    var category = document.getElementById("category1").value;
                    var amount = document.getElementById("amount1").value;
                } else {
                    var category = document.getElementById("category2").value;
                    var amount = document.getElementById("amount2").value;
                }
                var transaction = {
                    TransactionTime: getCurrentTime(),
                    Amount: amount * value,
                    Category: category
                };
                location.reload();
                var xhr = new XMLHttpRequest();
                xhr.open('POST', 'moneybuddy-production.up.railway.app/addtransaction', true);
                xhr.setRequestHeader('Content-Type', 'application/json');
                xhr.send(JSON.stringify(transaction));
                Get()
        }

        function getCurrentTime() {
            return new Date().toISOString();
        }


    // Get the modal
    var plusModal = document.getElementById("plusModal");    
    var minModal = document.getElementById("minModal");    
    // Get the button that opens the modal
    var open_plusModal = document.getElementById("open_plusModal");
    var open_minModal = document.getElementById("open_minModal");
    // Get the <span> element that closes the modal
    var close_plusModal = document.getElementById("close_plusModal");
    var close_minModal = document.getElementById("close_minModal");

    // When the user clicks the button, open the modal 
    open_plusModal.onclick = function() {
        plusModal.style.display = "block"; }
    open_minModal.onclick = function() {
        minModal.style.display = "block"; }

    // When the user clicks on <span> (x), close the modal
    close_plusModal.onclick = function() {
        plusModal.style.display = "none";}
    close_minModal.onclick = function() {
        minModal.style.display = "none";}
    // When the user clicks anywhere outside of the modal, close it
    window.onclick = function(event) {
        if (event.target == plusModal) {
        plusModal.style.display = "none";}
        if (event.target == minModal) {
            minModal.style.display = "none";}
    }
});