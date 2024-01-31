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
    
    // load header and footer in every file
    $(function(){
        $("#includedheader").load("support/header.html");
        $("#includedfooter").load("support/footer.html"); 
    });

    // function that makes text visible by click
    function toggleVisibility() {
        var hiddenText = document.getElementById('hiddenText');
        hiddenText.style.display = (hiddenText.style.display === 'none' || hiddenText.style.display === '') ? 'inline' : 'none';}
 
});