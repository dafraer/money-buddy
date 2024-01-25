document.addEventListener('DOMContentLoaded', function () {

     // api request

    /* 
        const apiUrl = '*'; // Kamil, replace * with go url

        const data = {

            category1: {
                name: "{{.Category1.Name}}",
                amount: {{.Category1.Amount}}
            },
            category2: {
                name: "{{.Category2.Name}}",
                amount: {{.Category2.Amount}}
            },
            category3: {
                name: "{{.Category3.Name}}",
                amount: {{.Category3.Amount}}
            },
            category4: {
                name: "{{.Category4.Name}}",
                amount: {{.Category4.Amount}}
            },
            category5: {
                name: "{{.Category5.Name}}",
                amount: {{.Category5.Amount}}
            },

            income: {{.Income}},
            
            expenditure: {{.Expenditure}}

            piggybank_target_amount: {{.PiggyBank.TargetAmount}}

            piggybank_target_date: {{.PiggyBank.TargetDate}}

            piggybank_target_date: {{.PiggyBank.Balance}}

            currency: {{.Currency}}
        };

        fetch(*, {              // Kamil, replace * with go url
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },

            body: JSON.stringify(data),
        })

        // check that everything is working

        .then(response => {      
            if (!response.ok) {
                throw new Error('Failed to update data');
            }
            console.log('Data updated successfully');
        })
        .catch(error => {
            console.error('Error updating data:', error);
        }); */

    const banana_image = document.querySelectorAll('.banana_image');
    const newAmount = document.getElementById('newAmount');
    // var targetAmount = {{.PiggyBank.TargetAmount}}

    
    // func that changes opacity
    function change_opacity() {

        var opacity_new = parseFloat(newAmount.value);
        var opacity_decimal = (opacity_new / targetAmount);
  
        banana_image[0].style.opacity = 1 - opacity_decimal;
        banana_image[1].style.opacity = opacity_decimal;
    }
  
    change_opacity(); 
    newAmount.addEventListener('change', change_opacity);

    

    // load header and footer in every file
    $(function(){
        $("#includedheader").load("support/header.html");
        $("#includedfooter").load("support/footer.html"); 
    });


    // function that makes text visible by click
    function toggleVisibility() {
        var hiddenText = document.getElementById('hiddenText');
        hiddenText.style.display = (hiddenText.style.display === 'none' || hiddenText.style.display === '') ? 'inline' : 'none';}

    // list of variables /* new categories should appear here as well */
    // var food_0 = 300;
    // var bills_0 = 300;
    // var transport_0 = 300;
    // var total = 900;

    // manipulations for right pie chart (with every variable) /* new categories should appear here as well */

    /*  var top1_1 = (({{.Category1.Name}}/ total) * 100) + 0; let top1_f = String(top1_1) + '%';
    var top2_1 = (({{.Category2.Name}}/ total) * 100) + top1_1; let top2_f = String(top2_1) + '%';
    var top3_1 = (({{.Category3.Name}}/ total) * 100) + top2_1; let top3_f = String(top3_1) + '%';
    var top4_1 = (({{.Category4.Name}}/ total) * 100) + top3_1; let top4_f = String(top4_1) + '%';
    var other_1 = ((other/ total) * 100) + top4_1; let other_f = String(other_1) + '%';
    */ 
    // func that changes css variables /* new categories should appear here as well */
    /* document.documentElement.style.setProperty('--top1CSS', top1_f);
    document.documentElement.style.setProperty('--top2CSS', top2_f);
    document.documentElement.style.setProperty('--top3CSS', top3_f);
    document.documentElement.style.setProperty('--top4CSS', top4_f);
    document.documentElement.style.setProperty('--otherCSS', other_f); */

    // pie always full (even wuthout expenses)
    var full = 100;
    let full1 = String(full) + '%';
    document.documentElement.style.setProperty('--fullCSS', full1);

    /* if ({{.Expenditure}} == null) {
        document.documentElement.style.setProperty('--fullCSS', full);}   */

    // banana math
    // var raw_banana = {{.Expenditure}};
    // raw_banana = 900 / 0.5; // total expenses
    // document.getElementById('bananamath').innerHTML = raw_banana;

    /* how categories are working?
    let see according to category food:
    food_0 = amount of $ spend on that particular category
    food_1 = this amount in percentage (food_0 / total) * 100% - pay attention, for work of the pie chart it's working a bit different for other categories
    food_f = f stands for final. food_f is a string version of food_1 with % sign. it's neccessary for pie chart to work
    food_c = id to make text bold and stylish :)
    -foodCSS = CSS variable; equal to food_f, but it's neccessary for pie chart to work
    thx guys ^)--> */

});