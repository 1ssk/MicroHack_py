document.querySelector('form').addEventListener('submit', function(event) {
    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;
    const confirmPassword = document.getElementById('confirm-password').value;

    if (!username || !password || !confirmPassword) {
      event.preventDefault();
      alert('Пожалуйста, заполните все поля.');
    } else if (password !== confirmPassword) {
      event.preventDefault();
      alert('Пароли не совпадают.');
    }
  });