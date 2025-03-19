document.addEventListener("DOMContentLoaded", function() {
    // Обработчик формы для создания курса
    document.getElementById('create-course-form').addEventListener('submit', function(e) {
        e.preventDefault();
        const title = document.getElementById('course-title').value;

        fetch('http://localhost:8000/courses/create', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ title: title }),
        })
        .then(response => {
            if (response.ok) {
                alert('Курс успешно создан');
                window.location.reload();
            } else {
                return response.json().then(data => {
                    throw new Error(data.detail || 'Ошибка при создании курса');
                });
            }
        })
        .catch(error => {
            console.error('Ошибка:', error);
            alert(error.message || 'Ошибка при создании курса');
        });
    });


    // Обработчик формы для создания урока
    document.getElementById('create-lesson-form').addEventListener('submit', function(e) {
        e.preventDefault();
        const title = document.getElementById('lesson-title').value;
        const url = document.getElementById('lesson-url').value;
        const courseId = document.getElementById('lesson-course').value;

        fetch('http://localhost:8000/lessons/create', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                title: title,
                url: url,
                course_id: courseId
            }),
        })
        .then(response => {
            if (response.ok) {
                alert('Урок успешно создан');
                window.location.reload();
            } else {
                return response.json().then(data => {
                    throw new Error(data.detail || 'Ошибка при создании урока');
                });
            }
        })
        .catch(error => {
            console.error('Ошибка:', error);
            alert(error.message || 'Ошибка при создании урока');
        });
    });

    
        // Обработчик формы для привязки курса к пользователю
    document.getElementById('assign-course-form').addEventListener('submit', function(e) {
        e.preventDefault();
        const userId = document.getElementById('user-select').value;
        const courseId = document.getElementById('course-select').value;

        fetch('http://localhost:8000/user/courses/create', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                user_id: userId,
                course_id: courseId
            }),
        })
        .then(response => {
            if (response.ok) {
                alert('Курс успешно привязан');
                window.location.reload();
            } else {
                return response.json().then(data => {
                    throw new Error(data.detail || 'Ошибка при привязке курса');
                });
            }
        })
        .catch(error => {
            console.error('Ошибка:', error);
            alert(error.message || 'Ошибка при привязке курса');
        });
    });
    
    // Обработчик кнопок удаления курса
    document.querySelectorAll('.delete-course-btn').forEach(button => {
        button.addEventListener('click', function() {
            const courseId = this.getAttribute('data-course-id');
            if (confirm(`Вы уверены, что хотите удалить курс с ID ${courseId} и все связанные уроки?`)) {
                fetch('/admin/delete-course', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/x-www-form-urlencoded',
                    },
                    body: `course_id=${encodeURIComponent(courseId)}`,
                })
                .then(response => {
                    if (response.ok) {
                        alert('Курс успешно удален');
                        window.location.reload();
                    } else {
                        alert('Ошибка при удалении курса');
                    }
                })
                .catch(error => {
                    console.error('Ошибка:', error);
                    alert('Ошибка при удалении курса');
                });
            }
        });
    });

    // Обработчик кнопок удаления пользователя
    document.querySelectorAll('.delete-user-btn').forEach(button => {
        button.addEventListener('click', function() {
            const userId = this.getAttribute('data-user-id');
            if (confirm(`Вы уверены, что хотите удалить пользователя с ID ${userId}?`)) {
                fetch('/admin/delete-user', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/x-www-form-urlencoded',
                    },
                    body: `user_id=${encodeURIComponent(userId)}`,
                })
                .then(response => {
                    if (response.ok) {
                        alert('Пользователь успешно удален');
                        window.location.reload();
                    } else {
                        alert('Ошибка при удалении пользователя');
                    }
                })
                .catch(error => {
                    console.error('Ошибка:', error);
                    alert('Ошибка при удалении пользователя');
                });
            }
        });
    });

    // Обработчик кнопок удаления урока
    document.querySelectorAll('.delete-lesson-btn').forEach(button => {
        button.addEventListener('click', function() {
            const lessonId = this.getAttribute('data-lesson-id');
            if (confirm(`Вы уверены, что хотите удалить урок с ID ${lessonId}?`)) {
                fetch('/admin/delete-lesson', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/x-www-form-urlencoded',
                    },
                    body: `lesson_id=${encodeURIComponent(lessonId)}`,
                })
                .then(response => {
                    if (response.ok) {
                        alert('Урок успешно удален');
                        window.location.reload();
                    } else {
                        alert('Ошибка при удалении урока');
                    }
                })
                .catch(error => {
                    console.error('Ошибка:', error);
                    alert('Ошибка при удалении урока');
                });
            }
        });
    });
});