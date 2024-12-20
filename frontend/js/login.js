document.getElementById('entry-btn').addEventListener('click', async () => {
    const login = document.getElementById('login').value;
    const password = document.getElementById('password').value;

    const requestBody = {
        "email": login,
        "password": password
    };

    try {
        const response = await fetch('http://localhost:5240/api/Account/Login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(requestBody)
        });

        console.log('Статус ответа:', response.status);

        if (!response.ok) {
            const errorText = await response.text(); // Получаем текст ошибки
            throw new Error(`Ошибка входа: ${response.status} ${errorText}`);
        }

        // Получаем текст ответа (токен)
        const token = await response.text();
        console.log('Ответ от сервера (текст):', token);

        // Сохраняем токен в localStorage
        localStorage.setItem('jwtToken', token);

        alert('Вход выполнен успешно!');

        // Здесь можно добавить логику для перехода на другую страницу или выполнения других действий
        window.location.href = '../html/subAfterAuth.html';
    } catch (error) {
        console.error('Ошибка:', error);
        alert('Не удалось войти. Проверьте логин и пароль.');
    }
});


// document.getElementById('reg-btn').addEventListener('click', async () => {
//     const name = document.getElementById('name').value;
//     const email = document.getElementById('email').value;
//     const password = document.getElementById('password').value;

//     const requestBody = {
//         "email": login,
//         "password": password
//     };

//     try {
//         const response = await fetch('http://localhost:5240/api/Account/Login', {
//             method: 'POST',
//             headers: {
//                 'Content-Type': 'application/json'
//             },
//             body: JSON.stringify(requestBody)
//         });

//         console.log('Статус ответа:', response.status);

//         if (!response.ok) {
//             const errorText = await response.text(); // Получаем текст ошибки
//             throw new Error(`Ошибка входа: ${response.status} ${errorText}`);
//         }

//         // Получаем текст ответа (токен)
//         const token = await response.text();
//         console.log('Ответ от сервера (текст):', token);

//         // Сохраняем токен в localStorage
//         localStorage.setItem('jwtToken', token);

//         alert('Вход выполнен успешно!');

//         // Здесь можно добавить логику для перехода на другую страницу или выполнения других действий
//         window.location.href = '../html/subAfterAuth.html';
//     } catch (error) {
//         console.error('Ошибка:', error);
//         alert('Не удалось войти. Проверьте логин и пароль.');
//     }
// });

// // Пример использования токена при нажатии на другую кнопку
// // document.getElementById('some-other-btn').addEventListener('click', async () => {
// //     const token = localStorage.getItem('jwtToken');

// //     if (!token) {
// //         alert('Пожалуйста, войдите в систему.');
// //         return;
// //     }

// //     // Используйте токен для выполнения защищенного запроса
// //     try {
// //         const response = await fetch('https://your-api-url.com/protected-resource', {
// //             method: 'GET',
// //             headers: {
// //                 'Authorization': 'Bearer ' + token
// //             }
// //         });

// //         if (!response.ok) {
// //             throw new Error('Ошибка доступа к ресурсу: ' + response.statusText);
// //         }

// //         const data = await response.json();
// //         console.log('Данные:', data);
// //     } catch (error) {
// //         console.error('Ошибка:', error);
// //     }
// // });