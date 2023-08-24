// modifyRedirect.js

document.addEventListener('DOMContentLoaded', function() {
  const penRows = document.querySelectorAll('.pen-row');

  penRows.forEach(row => {
    row.addEventListener('click', function() {
      const penID = this.id.replace('penRow', '');
      window.location.href = `/modify/${penID}`;
    });
  });
});
