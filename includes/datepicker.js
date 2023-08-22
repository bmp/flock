// includes/datepicker.js

// Get all input elements with type="date"
const dateInputs = document.querySelectorAll('input[type="date"]');

// Loop through the date inputs and set the date format to display only the year
dateInputs.forEach(input => {
  // Get the year from the date value
  const currentYear = new Date().getFullYear();
  const minYear = parseInt(input.min);
  const maxYear = parseInt(input.max);

  // Set the new value for the input
  input.value = currentYear;

  // Add a change event listener to update the value based on user input
  input.addEventListener('input', () => {
    const year = parseInt(input.value);
    if (year < minYear) {
      input.value = minYear;
    } else if (year > maxYear) {
      input.value = maxYear;
    }
  });
});
