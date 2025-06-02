// Simple test to check API connectivity
fetch('http://localhost:8090/cameras')
  .then(response => {
    console.log('Response status:', response.status);
    return response.json();
  })
  .then(data => {
    console.log('Camera data:', data);
  })
  .catch(error => {
    console.error('Error:', error);
  });
