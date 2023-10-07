BASE_URL = "http://localhost:9000/api/v1";


function viewProduct(e) {
    let card = e.target.parentElement.parentElement;
    let img = card.querySelector('img').cloneNode(true);
    let title = $(card).find(".card-title").text()
    let description = $(card).find(".card-text").text();
    let modalTitle = $('#productDetail').find('.modal-title');
    let modalText = $('#productDetail').find('.modal-text');
    let modalImage = $('#productDetail').find('.modal-img');
    let nonProfitMatch = $('#productDetail').find('.non-profit-matched');

    score(description).then(name => {
        nonProfitMatch.text(name);
    });

    modalImage.empty();
    modalImage.append(img);
    modalTitle.text(title);
    modalText.text(description);

    $('#productDetail').modal('toggle');
}

async function score(input) {
    let scoringEndpoint = `${BASE_URL}/score/${input}`;
    const response = await fetch(scoringEndpoint);
    const match = await response.json();
    if (match != null && match.name != null && match.name.length > 0) {
        return match.name;
    }
    return null;
}

async function matchNonProfit(subcategory) {
    let scoringEndpoint = `${BASE_URL}/match?subcategory=${subcategory}`;
    const response = await fetch(scoringEndpoint);
    const nonProfits = await response.json();
    if (nonProfits.length != null) {
        return nonProfits.Name;
    }
    return null;
}