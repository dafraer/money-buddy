document.addEventListener('DOMContentLoaded', function () {

    // api request function
    xhr = new XMLHttpRequest();
    xhr.open('GET', 'http://localhost:8000/getuserdata', true);
    xhr.onload = function () {
        if (this.status == 200) {
            var u = JSON.parse(this.responseText);
            u.Analytics.Income == null ? document.getElementById("totalInc").innerHTML = 0 : document.getElementById("totalInc").innerHTML = u.Analytics.Income;
            u.Analytics.Expenditure == null ? document.getElementById("totalSpent").innerHTML = 0 : document.getElementById("totalSpent").innerHTML = u.Analytics.Expenditure;
            for (let i = 0; i < 4; i++) {
                document.getElementById("name" + i).innerHTML = u.Analytics.Categories[i].Name;
                u.Analytics.Categories[i].Amount == null ? document.getElementById("amount" + i).innerHTML = 0 : document.getElementById("amount" + i).innerHTML = u.Analytics.Categories[i].Amount;
                // Hide category if amount is zero
                if (u.Analytics.Categories[i].Amount == 0) {
                    document.getElementById("category" + i).style.display = "none";
                }
            }
            u.Analytics.Categories[4].Amount == null ? document.getElementById("amountother").innerHTML = 0 : document.getElementById("amountother").innerHTML = u.Analytics.Categories[4].Amount;

            // Hide 'Other' category if amount is zero
            if (u.Analytics.Categories[4].Amount == 0) {
                document.getElementById("other").style.display = "none";
                document.getElementById("amountother").style.display = "none";
            }

            // List of variables and manipulations
            var amount0 = parseFloat(document.getElementById("amount0").textContent);
            var amount1 = parseFloat(document.getElementById("amount1").textContent);
            var amount2 = parseFloat(document.getElementById("amount2").textContent);
            var amount3 = parseFloat(document.getElementById("amount3").textContent);
            var other = parseFloat(document.getElementById("amountother").textContent);
            var totalSpent = parseFloat(document.getElementById("totalSpent").textContent);
            
            var other_f = ((other / totalSpent) * 100) + '%';
            var amount0_f = ((amount0 / totalSpent) * 100) + 0 + '%';
            var amount1_f = ((amount1 / totalSpent) * 100) + ((amount0 / totalSpent) * 100) + '%';
            var amount2_f = ((amount2 / totalSpent) * 100) + ((amount1 / totalSpent) * 100) + ((amount0 / totalSpent) * 100) + '%';
            var amount3_f = ((amount3 / totalSpent) * 100) + ((amount2 / totalSpent) * 100) + ((amount1 / totalSpent) * 100) + ((amount0 / totalSpent) * 100) + '%';

            // pie always full (even wuthout expenses)
            var full = 0;
            if (totalSpent == '0') {
                full = 100 + '%';
                document.documentElement.style.setProperty('--fullCSS', full);
                document.getElementById("noexpenses").innerHTML = "pie chart will change when you'll start tracking expenses!";
                // noexpenses
             } 
            else {
                document.documentElement.style.setProperty('--otherCSS', other_f);
                document.documentElement.style.setProperty('--top0CSS', amount0_f);
                document.documentElement.style.setProperty('--top1CSS', amount1_f);
                document.documentElement.style.setProperty('--top2CSS', amount2_f);
                document.documentElement.style.setProperty('--top3CSS', amount3_f);
            }
            // banana math
            bananamath = Math.round(totalSpent / 0.6);
            document.getElementById('bananamath').innerHTML = bananamath;
        }
    }
    xhr.send();
}); 
