let currentPage = 1;


function getToken() {
    const match = document.cookie.match(/token=([^;]+)/);
    return match ? match[1] : null;
}

function logout() {
    document.cookie = 'token=; path=/; expires=Thu, 01 Jan 1970 00:00:00 UTC';
    window.location.href = '/';
}

async function updateListings() {
    const filter = document.getElementById('filter').value;
    const res = await fetch(`/api/listings?page=${currentPage}&filter=${filter}`, {
        headers: {'Authorization': `Bearer ${getToken()}`}
    });
    const data = await res.json();
    const tbody = document.getElementById('listings');
    tbody.innerHTML = '';
    data.forEach(item => {
        const tr = document.createElement('tr');
        tr.innerHTML = `
      <td><input value="${item.title}" onchange="updateField(${item.id}, 'title', this.value)"></td>
      <td><input type="number" value="${item.price}" onchange="updateField(${item.id}, 'price', this.value)"></td>
      <td><input value="${item.city}" onchange="updateField(${item.id}, 'city', this.value)"></td>
      <td><button onclick="deleteListing(${item.id})">Удалить</button></td>
    `;
        tbody.appendChild(tr);
    });
    document.getElementById('page-info').textContent = currentPage;
}

async function deleteListing(id) {
    await fetch(`/api/listings/${id}`, {
        method: 'DELETE',
        headers: {'Authorization': `Bearer ${getToken()}`}
    });
    await updateListings();
    await updateAnalytics();
}

async function updateField(id, field, value) {
    await fetch(`/api/listings/${id}`, {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${getToken()}`
        },
        body: JSON.stringify({[field]: value})
    });
    await updateAnalytics();
}

async function updateAnalytics() {
    const res = await fetch('/api/analytics', {
        headers: {'Authorization': `Bearer ${getToken()}`}
    });
    const data = await res.json();
    document.getElementById('analytics').innerHTML = `
    Объявлений: ${data.total}<br>
    Средняя цена: ${data.avg_price} ₽<br>
    Топ-город: ${data.top_city}
  `;
}

function nextPage() {
    currentPage++;
    updateListings();
}

function prevPage() {
    if (currentPage > 1) {
        currentPage--;
        updateListings();
    }
}

if (document.getElementById('add-form')) {
    document.getElementById('add-form').addEventListener('submit', async e => {
        e.preventDefault();
        const title = document.getElementById('title').value;
        const price = parseFloat(document.getElementById('price').value);
        const city = document.getElementById('city').value;

        await fetch('/api/listings', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${getToken()}`
            },
            body: JSON.stringify({title, price, city})
        });

        currentPage = 1;
         updateListings();
         updateAnalytics();
        e.target.reset();
    });

    updateListings();
    updateAnalytics();
}


