let currentPage = 1;
let editingId = null;
let hasNextPage = true;
let hasPrevPage = false;


function formatDate(isoDateString) {
    const date = new Date(isoDateString);
    return date.toLocaleString("ru-RU", {
        timeZone: "UTC",
        day: "2-digit",
        month: "2-digit",
        year: "numeric",
        hour: "2-digit",
        minute: "2-digit",
    });
}

function logout() {
    fetch("/api/logout", {
        method: "POST",
    }).then(res => {
        if (res.redirected) {
            window.location.href = res.url;
        }
    });
}

async function updateListings() {
    const filter = document.getElementById('filter').value;
    const res = await fetch(`/api/listings?page=${currentPage}&filter=${filter}`, {
    });
    const data = await res.json();
    const tbody = document.getElementById('listings');
    tbody.innerHTML = '';

    if (data === null) {
        hasNextPage = false;
        if (currentPage > 1) {
            currentPage--;
             await updateListings();
            return;
        }
        const tr = document.createElement('tr');
        tr.innerHTML = `<td colspan="9" style="text-align:center; padding: 10px;">Нет объявлений</td>`;
        tbody.appendChild(tr);
    } else {
        hasNextPage = data.length === 10;
        hasPrevPage = currentPage > 1;

        data.forEach((item, index) => {

            const tr = document.createElement('tr');
            const itemNumber = (currentPage - 1) * 10 + index + 1; // глобальный номер
            tr.innerHTML = `
       <td>${itemNumber}</td>
      <td><strong>${item.Name}</strong></td>
      <td>${item.Typel}</td>
      <td>${item.Description}</td>
      <td>${item.Status}</td>
      <td>${item.Price}</td>
      <td>${item.City}</td>
      <td>${formatDate(item.Date_created)}</td>
      
      <td><button class="action-button edit-btn" onclick=openModal(${JSON.stringify(item)})>Изменить</button> <button class="action-button delete-btn" onclick="deleteListing(${item.ID})">Удалить</button></td>
    `;
            //onclick=openModal(${JSON.stringify(item)}))
            tbody.appendChild(tr);
        });
    }
    document.getElementById('page-info').textContent = currentPage;
    // отключим/включим кнопки
    document.getElementById("prev-button").disabled = !hasPrevPage;
    document.getElementById("next-button").disabled = !hasNextPage;
}

async function deleteListing(id) {
    await fetch(`/api/listings/${id}`, {
        method: 'DELETE',
    });
    await updateListings();
    await updateAnalytics();
}

function openModal(listing = null) {
    // document.getElementById('modal').classList.remove('hidden');
    const saveBtn = document.getElementById('modal-save-button');
    const titleInput = document.getElementById('modal-title');
    const priceInput = document.getElementById('modal-price');
    const cityInput = document.getElementById('modal-city');
    const typeInput = document.getElementById('modal-type');
    const descriptionInput = document.getElementById('modal-description');
    const statusInput = document.getElementById('modal-status');
    const label = document.getElementById('modal-title-label')

    if (listing) {
        // режим редактирования
        editingId = listing.ID;
        typeInput.value = listing.Typel
        descriptionInput.value = listing.Description;
        statusInput.value = listing.Status;
        titleInput.value = listing.Name;
        priceInput.value = listing.Price;
        cityInput.value = listing.City;
        saveBtn.textContent = "Сохранить";
        label.textContent = "Редактировать объявление";
    } else {
        // режим добавления
        editingId = null;
        titleInput.value = "";
        priceInput.value = "";
        cityInput.value = "";
        saveBtn.textContent = "Создать";
        label.textContent = "Новое объявление";
    }
    saveBtn.onclick = () => {
        if (editingId) {
            updateListing(editingId);
        } else {
            createListing();
        }
    }
    document.getElementById('modal').classList.remove('hidden');
}

async function updateAnalytics() {
    const res = await fetch('/api/analytics', {
        // headers: {'Authorization': `Bearer ${getToken()}`}
    });
    const data = await res.json();
    document.getElementById('analytics').innerHTML = `
    Объявлений: ${data.total}<br>
    Средняя цена: ${data.avg_price} ₽<br>
    Топ-город: ${data.top_city}
  `;
}

function nextPage() {
    // currentPage++;
    // updateListings();
    if (!hasNextPage) return;
    currentPage++;
    updateListings();
}

function prevPage() {
    // if (currentPage > 1) {
    //     currentPage--;
    //     updateListings();
    // }
    if (currentPage <= 1) return;
    currentPage--;
    updateListings();
}

function closeModal() {
    document.getElementById('modal').classList.add('hidden');
    document.getElementById('modal-title').value = '';
    document.getElementById('modal-price').value = '';
    document.getElementById('modal-city').value = '';
    document.getElementById('modal-type').value = '';
    document.getElementById('modal-description').value = '';
    document.getElementById('modal-status').value = '';
}

async function createListing() {
    const title = document.getElementById('modal-title').value;
    const type = document.getElementById('modal-type').value;
    const description = document.getElementById('modal-description').value;
    const status = document.getElementById('modal-status').value;
    const price = parseInt(document.getElementById('modal-price').value);
    const city = document.getElementById('modal-city').value;

    await fetch('/api/listings', {
        method: 'POST',
        // headers: {
        //     'Content-Type': 'application/json',
        //     'Authorization': `Bearer ${getToken()}`
        // },
        body: JSON.stringify({ title, type, description, status, price, city })
    }).then(res => {
        if (res.ok){
            alert("успешно добавленно")
        }
        }
    );

    closeModal();
    currentPage = 1;
    updateListings();
    updateAnalytics();
}

async function loadCities() {
    const res = await fetch("/api/cities", {
        // headers: { Authorization: `Bearer ${getToken()}` }
    });
    const cities = await res.json();
    if (cities===null) return;
    const select = document.getElementById("filter");
    cities.forEach(city => {
        const option = document.createElement("option");
        option.value = city;
        option.textContent = city;
        select.appendChild(option);
    });
}

async function updateListing(id) {
    const title = document.getElementById('modal-title').value;
    const type = document.getElementById('modal-type').value;
    const description = document.getElementById('modal-description').value;
    const status = document.getElementById('modal-status').value;
    const price = parseInt(document.getElementById('modal-price').value);
    const city = document.getElementById('modal-city').value;

    await fetch(`/api/listings/${id}`, {
        method: 'PUT',
        // headers: {
        //     'Content-Type': 'application/json',
        //     'Authorization': `Bearer ${getToken()}`
        // },
        body: JSON.stringify({ title, type, description, status, price, city })
    }).then(res => {
            if (res.ok){
                alert("успешно Изменено")
            }
        }
    );

    closeModal();
    updateListings();
    updateAnalytics();
}
