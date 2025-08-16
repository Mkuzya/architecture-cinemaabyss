const fs = require('fs');

// Читаем файл
let content = fs.readFileSync('CinemaAbyss.postman_collection.json', 'utf8');

// Исправляем все тесты ID matches or exists
const idTestPattern = /pm\.test\("ID matches or exists", function \(\) \{[\s\S]*?\}\);/g;

const replacement = `pm.test("ID matches or exists", function () {
    var jsonData = pm.response.json();
    var userId = pm.collectionVariables.get("userId");
    var movieId = pm.collectionVariables.get("movieId");
    var paymentId = pm.collectionVariables.get("paymentId");
    var subscriptionId = pm.collectionVariables.get("subscriptionId");
    var msMovieId = pm.collectionVariables.get("msMovieId");
    
    // Проверяем, какая переменная должна быть установлена
    var expectedId = null;
    if (userId && userId !== "" && userId !== "{{userId}}") {
        expectedId = parseInt(userId);
    } else if (movieId && movieId !== "" && movieId !== "{{movieId}}") {
        expectedId = parseInt(movieId);
    } else if (paymentId && paymentId !== "" && paymentId !== "{{paymentId}}") {
        expectedId = parseInt(paymentId);
    } else if (subscriptionId && subscriptionId !== "" && subscriptionId !== "{{subscriptionId}}") {
        expectedId = parseInt(subscriptionId);
    } else if (msMovieId && msMovieId !== "" && msMovieId !== "{{msMovieId}}") {
        expectedId = parseInt(msMovieId);
    }
    
    if (expectedId) {
        pm.expect(jsonData.id).to.equal(expectedId);
    } else {
        pm.expect(jsonData.id).to.exist;
    }
});`;

content = content.replace(idTestPattern, replacement);

// Записываем исправленный файл
fs.writeFileSync('CinemaAbyss.postman_collection.json', content);

console.log('ID tests fixed!');
