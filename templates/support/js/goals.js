document.addEventListener('DOMContentLoaded', function () {    

    var user = {};
    Get();
    document.getElementById("Form").addEventListener("submit", Update);

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

                // Set opacity of green banana image based on balance percentage
                var greenBananaOpacity = balancePercentage >= 100 ? '0' : (1 - (balancePercentage / 100)).toString();
                document.getElementById("image1").style.opacity = greenBananaOpacity;

                // Set opacity of yellow banana image based on balance percentage
                var yellowBananaOpacity = balancePercentage >= 100 ? '1' : (balancePercentage / 100).toString();
                document.getElementById("image2").style.opacity = yellowBananaOpacity;

                // Update progress label with rounded percentage
                document.getElementById("progress-label").innerText = "Your progress: " + balancePercentage + "%";
            }
        };
        xhr.send();
    }

    function Update(e) {
        e.preventDefault();
        var amount = parseFloat(document.getElementById("amount").value); 
        var newAmount = parseFloat(document.getElementById("newAmount").value);
        var newDate = document.getElementById("newDate").value;
        if (isNaN(amount)) {
            amount = 0;
        }
        if (isNaN(newAmount)) {
            newAmount = user.PiggyBank.TargetAmount;
        }
        if (newDate == "") {
            newDate = user.PiggyBank.TargetDate;
        }
        var PiggyBank = {
            TargetAmount: newAmount,
            TargetDate: newDate,
            Balance: amount
        };

        var xhr = new XMLHttpRequest();
        xhr.open('POST', 'http://localhost:8000/postpiggybank', true);
        xhr.setRequestHeader('Content-Type', 'application/json');
        xhr.send(JSON.stringify(PiggyBank));
        xhr.onload = function() {
            Get();
        }
    }
});
