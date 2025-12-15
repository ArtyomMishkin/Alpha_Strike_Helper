let currentPage = 1;
let pageSize = 9999; // ИСПРАВЛЕНО: было 20, теперь загружает ВСЕ карты сразу
let allCards = [];
let filteredCards = [];

const API_BASE_URL = 'http://localhost:8080/api/v1';

document.addEventListener('DOMContentLoaded', function() {
    console.log('DOM loaded, initializing cards...');
    loadCards();
    setupFilterHandlers();
    setupActionHandlers();
});

async function loadCards(filters = {}) {
    try {
        console.log('Loading cards with filters:', filters);
        
        const params = new URLSearchParams();
        params.set('page', currentPage);
        params.set('pagesize', pageSize); // 9999 - загружает всё
        
        if (filters.role && filters.role !== '') params.append('role', filters.role);
        if (filters.size && filters.size !== '') params.append('size', filters.size);
        if (filters.faction && filters.faction !== '') params.append('faction', filters.faction);
        if (filters.type && filters.type !== '') params.append('type', filters.type);
        if (filters.techbase && filters.techbase !== '') params.append('techbase', filters.techbase);
        if (filters.name && filters.name !== '') params.append('name', filters.name);
        if (filters.pvmin && filters.pvmin !== '') params.append('pvmin', filters.pvmin);
        if (filters.pvmax && filters.pvmax !== '') params.append('pvmax', filters.pvmax);
        
        const url = `${API_BASE_URL}/cards?${params.toString()}`;
        console.log('Requesting:', url);
        
        const response = await fetch(url);
        
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        const result = await response.json();
        console.log('Response:', result);
        
        if (!result.data) {
            console.error('Invalid response structure:', result);
            allCards = [];
        } else {
            allCards = result.data;
        }
        
        console.log(`Loaded ${allCards.length} cards`);
        displayCards(allCards);
        updatePaginationInfo(result.total || 0, result.totalpages || 0);
        
    } catch (error) {
        console.error('Error loading cards:', error);
        
        // Fallback to demo data
        console.log('Using fallback demo data...');
        loadDemoCards();
    }
}

// DEMO FALLBACK
function loadDemoCards() {
    allCards = [
        {
            id: 1,
            name: 'Commando',
            modelnumber: 'COM-2D',
            type: 'Medium Mech',
            size: 'Medium',
            faction: 'Inner Sphere',
            role: 'Skirmisher',
            techbase: 'Clan',
            pointvalue: 35,
            move: '8/12',
            tmm: 0,
            armor: 8,
            structure: 8,
            damageshort: 3,
            damagemedium: 2,
            damagelong: 1,
            overheat: 2,
            abilities: 'Jump Jets'
        },
        {
            id: 2,
            name: 'Locust',
            modelnumber: 'LCT-1V',
            type: 'Light Mech',
            size: 'Light',
            faction: 'Inner Sphere',
            role: 'Scout',
            techbase: 'IS',
            pointvalue: 18,
            move: '16/24',
            tmm: 1,
            armor: 3,
            structure: 2,
            damageshort: 0,
            damagemedium: 0,
            damagelong: 0,
            overheat: 0,
            abilities: 'None'
        }
    ];
    
    displayCards(allCards);
}

function displayCards(cards) {
    console.log('Displaying', cards.length, 'cards');
    
    const container = document.getElementById('cards-container');
    if (!container) {
        console.error('Cards container not found!');
        return;
    }
    
    if (!cards || cards.length === 0) {
        container.innerHTML = '<div class="empty-state"><h3>No cards found</h3><p>Try adjusting your filters</p></div>';
        return;
    }
    
    container.innerHTML = cards.map(card => createCardHTML(card)).join('');
    addCardEventListeners();
}

// HTML GENERATION
function createCardHTML(card) {
    return `
        <div class="card card-glow" data-card-id="${card.id}">
            <div class="card-image">
                <div style="background: linear-gradient(135deg, #457b9d 0%, #1d3557 100%); height: 100px; display: flex; align-items: center; justify-content: center; color: white; font-weight: bold;">
                    ${card.size || 'Unknown'}
                </div>
                <div class="card-type-badge">${card.type || 'BattleMech'}</div>
            </div>
            
            <div class="card-content">
                <div class="card-header">
                    <div>
                        <h3 class="card-title">${escapeHtml(card.name || 'Unknown')}</h3>
                        <p class="card-model">${escapeHtml(card.modelnumber || '')}</p>
                    </div>
                    <div class="card-pv">${card.pointvalue || 0} PV</div>
                </div>
                
                <div class="card-stats">
                    <div class="stat-item">
                        <div class="stat-label">Move</div>
                        <div class="stat-value">${escapeHtml(card.move || '-')}</div>
                    </div>
                    <div class="stat-item">
                        <div class="stat-label">TMM</div>
                        <div class="stat-value">${card.tmm || 0}</div>
                    </div>
                    <div class="stat-item">
                        <div class="stat-label">Armor</div>
                        <div class="stat-value">${card.armor || 0}</div>
                    </div>
                    <div class="stat-item">
                        <div class="stat-label">Struct</div>
                        <div class="stat-value">${card.structure || 0}</div>
                    </div>
                </div>
                
                <div style="font-size: 0.85rem; color: #666; margin-top: 0.5rem;">
                    <div>Role: <strong>${escapeHtml(card.role || '-')}</strong></div>
                    <div>Faction: <strong>${escapeHtml(card.faction || '-')}</strong></div>
                </div>
                
                <div class="card-actions">
                    <button class="card-btn card-btn-view" onclick="viewCard(${card.id})">View</button>
                    <button class="card-btn card-btn-add add-to-lance" data-card-id="${card.id}">Lance</button>
                    <button class="card-btn card-btn-add add-to-star" data-card-id="${card.id}">Star</button>
                </div>
            </div>
        </div>
    `;
}

// ESCAPE HTML
function escapeHtml(text) {
    if (!text) return '';
    
    const map = {
        '&': '&amp;',
        '<': '&lt;',
        '>': '&gt;',
        '"': '&quot;',
        "'": '&#039;'
    };
    
    return String(text).replace(/[&<>"']/g, m => map[m]);
}

function setupFilterHandlers() {
    const applyBtn = document.getElementById('apply-filters-btn');
    if (!applyBtn) {
        console.warn('Apply filters button not found');
        return;
    }
    
    applyBtn.addEventListener('click', function(e) {
        e.preventDefault();
        console.log('Apply filters clicked');
        
        const filters = {
            role: document.getElementById('role-filter')?.value || '',
            size: document.getElementById('size-filter')?.value || '',
            faction: document.getElementById('faction-filter')?.value || '',
            type: document.getElementById('type-filter')?.value || '',
            techbase: document.getElementById('techbase-filter')?.value || '',
            name: document.getElementById('name-filter')?.value || '',
            pvmin: document.getElementById('pvmin-filter')?.value || '',
            pvmax: document.getElementById('pvmax-filter')?.value || ''
        };
        
        console.log('Filters:', filters);
        currentPage = 1;
        loadCards(filters);
    });
}

function addCardEventListeners() {
    // Add to Lance
    document.querySelectorAll('.add-to-lance').forEach(btn => {
        btn.addEventListener('click', function(e) {
            e.stopPropagation();
            const cardId = this.getAttribute('data-card-id');
            addToLance(cardId);
        });
    });
    
    // Add to Star
    document.querySelectorAll('.add-to-star').forEach(btn => {
        btn.addEventListener('click', function(e) {
            e.stopPropagation();
            const cardId = this.getAttribute('data-card-id');
            addToStar(cardId);
        });
    });
}

function setupActionHandlers() {
    // Setup
    console.log('Action handlers setup ready');
}

function viewCard(cardId) {
    const card = allCards.find(c => c.id === parseInt(cardId));
    if (card) {
        console.log('Viewing card:', card);
        alert(`Card: ${card.name}\n${card.type}\n${card.pointvalue} PV`);
    }
}

function addToLance(cardId) {
    console.log('Adding to Lance:', cardId);
    const lances = JSON.parse(localStorage.getItem('lances') || '[]');
    if (!lances.includes(parseInt(cardId))) {
        lances.push(parseInt(cardId));
        localStorage.setItem('lances', JSON.stringify(lances));
        console.log('Card added to lance');
        alert('Card added to lance!');
    } else {
        alert('Card already in lance');
    }
}

function addToStar(cardId) {
    console.log('Adding to Star:', cardId);
    const stars = JSON.parse(localStorage.getItem('stars') || '[]');
    if (!stars.includes(parseInt(cardId))) {
        stars.push(parseInt(cardId));
        localStorage.setItem('stars', JSON.stringify(stars));
        console.log('Card added to star');
        alert('Card added to star!');
    } else {
        alert('Card already in star');
    }
}

function updatePaginationInfo(total, totalPages) {
    const pageInfo = document.getElementById('page-info');
    if (pageInfo) {
        pageInfo.textContent = `Page ${currentPage} of ${totalPages} (total: ${total})`;
    }
}

function nextPage() {
    currentPage++;
    loadCards();
    window.scrollTo({top: 0, behavior: 'smooth'});
}

function prevPage() {
    if (currentPage > 1) {
        currentPage--;
        loadCards();
        window.scrollTo({top: 0, behavior: 'smooth'});
    }
}

function searchCards(query) {
    console.log('Searching:', query);
    if (!query || !query.trim()) {
        loadCards();
        return;
    }
    
    const searchFilters = {
        name: query,
    };
    currentPage = 1;
    loadCards(searchFilters);
    window.scrollTo({top: 0, behavior: 'smooth'});
}

// Expose to global scope
window.addToLance = addToLance;
window.addToStar = addToStar;
window.viewCard = viewCard;
window.searchCards = searchCards;
window.nextPage = nextPage;
window.prevPage = prevPage;
window.loadCards = loadCards;