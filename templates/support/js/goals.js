document.addEventListener('DOMContentLoaded', function () {    
    var user = {};
    var timern = new Date(); // Get the current date
    Get();
    document.getElementById("Form").addEventListener("submit", Update);

    // ajax request func
    function Get() {
        var xhr = new XMLHttpRequest();
        xhr.open('GET', 'http://localhost:8000/getuserdata', true);
        xhr.onload = function () {
            if (this.status == 200) {
                var u = JSON.parse(this.responseText);
                user = u;
                document.getElementById("targetAmount").innerHTML = u.PiggyBank.TargetAmount;
                document.getElementById("targetDate").innerHTML = u.PiggyBank.TargetDate;
                document.getElementById("balance").innerHTML = u.PiggyBank.Balance;

                // Calculate balance percentage rounded to the nearest integer
                var balancePercentage = Math.round((u.PiggyBank.Balance / u.PiggyBank.TargetAmount) * 100);

                if (isNaN(balancePercentage)) {
                    balancePercentage = 0;
                }
                if (!isFinite(balancePercentage)) {
                    balancePercentage = 0;
                }

                // Set opacity of green banana image based on balance percentage
                var greenBananaOpacity = balancePercentage >= 100 ? '0' : (1 - (balancePercentage / 100)).toString();
                document.getElementById("image1").style.opacity = greenBananaOpacity;

                // Set opacity of yellow banana image based on balance percentage
                var yellowBananaOpacity = balancePercentage >= 100 ? '1' : (balancePercentage / 100).toString();
                document.getElementById("image2").style.opacity = yellowBananaOpacity;

                // Update progress label with rounded percentage
                document.getElementById("progress-label").innerText = "Your progress: " + balancePercentage + "%";
                
                // calculate how many dollares left
                var targetDate = new Date(u.PiggyBank.TargetDate);
                var total_seconds = Math.abs(targetDate - timern) / 1000;  
                var days_difference = Math.trunc(total_seconds / (60 * 60 * 24)) + 1;
                console.log(days_difference);
                var targetAmount = u.PiggyBank.TargetAmount;
                var balance = u.PiggyBank.Balance;

                // check that it's not 0 or 1 days letf
                if ((days_difference == 0) || (days_difference == 1)) {
                    var dolares = targetAmount;
                    document.getElementById("dolares").innerHTML = dolares;
                }

                // check if user reached target amount
                else if (balance==targetAmount) {
                    var success = "You've reached your goal!";
                    document.getElementById("success").innerHTML = success;
                    document.getElementById("daily").style.display = "none";
                } 

                else {
                    // how many dollars per day are needed = dolares
                    var dolares = (targetAmount - balance) / days_difference;
                    dolares = dolares.toFixed(2);

                    // check that dolares is not Nan or Infinity
                    if (isNaN(dolares)) {
                        dolares = 0;
                    }
                    
                    if (!isFinite(dolares)) {
                        dolares = 0;
                    }

                    document.getElementById("dolares").innerHTML = dolares;
                }

            }
        };
        xhr.send();
    }

    // update database with new variables
    function Update(e) {
        e.preventDefault();
        var amount = parseFloat(document.getElementById("amount").value); 
        var newAmount = parseFloat(document.getElementById("newAmount").value);
        var newDate = document.getElementById("newDate").value;
        
        if (isNaN(amount)) {
            amount = 0;
        }

        if (amount < 0) {
            amount = 0;
            alert('Please add positive amount!');
        }

        if (isNaN(newAmount)) {
            newAmount = user.PiggyBank.TargetAmount;
        }

        if (newAmount < 0) {
            newAmount = user.PiggyBank.TargetAmount;
            alert('Please add positive amount!');
        }

        if (newDate == "") {
            newDate = user.PiggyBank.TargetDate;

        } else if (new Date(newDate) < timern) {
            alert('Please enter a future date!');
            newDate = user.PiggyBank.TargetDate;
        }      
        
        var PiggyBank = {
            TargetAmount: newAmount,
            TargetDate: newDate,
            Balance: amount
        };

        location.reload();

        var xhr = new XMLHttpRequest();
        xhr.open('POST', 'http://localhost:8000/postpiggybank', true);
        xhr.setRequestHeader('Content-Type', 'application/json');
        xhr.send(JSON.stringify(PiggyBank));
        xhr.onload = function() {
            Get();
            // Clear input fields after form submission
            document.getElementById("amount").value = "";
            document.getElementById("newAmount").value = "";
            document.getElementById("newDate").value = "";
        }
    }
});
