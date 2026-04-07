package server

import "net/http"

func (s *Server) dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(dashHTML))
}

const dashHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width,initial-scale=1.0">
<title>Prairie</title>
<link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500;700&display=swap" rel="stylesheet">
<style>
:root{--bg:#1a1410;--bg2:#241e18;--bg3:#2e261e;--rust:#e8753a;--leather:#a0845c;--cream:#f0e6d3;--cd:#bfb5a3;--cm:#7a7060;--gold:#d4a843;--green:#4a9e5c;--red:#c94444;--orange:#d4843a;--blue:#5b8dd9;--purple:#9d6bb8;--mono:'JetBrains Mono',monospace}
*{margin:0;padding:0;box-sizing:border-box}
body{background:var(--bg);color:var(--cream);font-family:var(--mono);line-height:1.5;font-size:13px;height:100vh;overflow:hidden;display:flex;flex-direction:column}
.hdr{padding:.7rem 1.2rem;border-bottom:1px solid var(--bg3);display:flex;justify-content:space-between;align-items:center;gap:1rem;flex-wrap:wrap;flex-shrink:0}
.hdr h1{font-size:.85rem;letter-spacing:2px}
.hdr h1 span{color:var(--rust)}
.hdr-left{display:flex;align-items:center;gap:1rem;flex-wrap:wrap}
.board-picker{padding:.4rem .6rem;background:var(--bg2);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem;min-width:180px}
.board-picker:focus{outline:none;border-color:var(--leather)}
.hdr-actions{display:flex;gap:.4rem}
.btn{font-family:var(--mono);font-size:.6rem;padding:.3rem .55rem;cursor:pointer;border:1px solid var(--bg3);background:var(--bg);color:var(--cd);transition:.15s}
.btn:hover{border-color:var(--leather);color:var(--cream)}
.btn-p{background:var(--rust);border-color:var(--rust);color:#fff}
.btn-p:hover{opacity:.85;color:#fff}
.btn-sm{font-size:.55rem;padding:.2rem .4rem}
.btn-del{color:var(--red);border-color:#3a1a1a}
.btn-del:hover{border-color:var(--red);color:var(--red)}

.board-area{flex:1;overflow-x:auto;overflow-y:hidden;padding:1rem;display:flex;gap:.8rem;align-items:flex-start;background:var(--bg)}
.column{background:var(--bg2);border:1px solid var(--bg3);min-width:260px;max-width:280px;display:flex;flex-direction:column;max-height:calc(100vh - 110px);flex-shrink:0}
.column.drag-over{border-color:var(--rust)}
.col-hdr{padding:.7rem .8rem;border-bottom:1px solid var(--bg3);display:flex;justify-content:space-between;align-items:center;gap:.4rem}
.col-name{font-family:var(--mono);font-size:.6rem;color:var(--leather);text-transform:uppercase;letter-spacing:1px;font-weight:700}
.col-count{font-family:var(--mono);font-size:.55rem;color:var(--cm);background:var(--bg);padding:.1rem .35rem}
.col-cards{padding:.5rem;overflow-y:auto;flex:1;min-height:60px;display:flex;flex-direction:column;gap:.4rem}
.col-add{padding:.4rem .8rem;border-top:1px solid var(--bg3)}
.col-add input{width:100%;padding:.35rem .5rem;background:var(--bg);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.65rem}
.col-add input:focus{outline:none;border-color:var(--leather)}

.card{background:var(--bg);border:1px solid var(--bg3);padding:.6rem .7rem;cursor:grab;display:flex;flex-direction:column;gap:.3rem;transition:border-color .15s}
.card:hover{border-color:var(--leather)}
.card.dragging{opacity:.4;cursor:grabbing}
.card-title{font-family:var(--mono);font-size:.7rem;color:var(--cream);font-weight:500;line-height:1.4}
.card-meta{display:flex;gap:.4rem;flex-wrap:wrap;align-items:center;font-size:.5rem;color:var(--cm)}
.card-assignee{background:var(--bg3);color:var(--cd);padding:.05rem .3rem;font-family:var(--mono);font-size:.5rem;text-transform:uppercase;letter-spacing:.5px}
.card-label{padding:.05rem .3rem;font-family:var(--mono);font-size:.5rem;text-transform:uppercase;letter-spacing:.5px;background:var(--rust);color:#fff}
.card-extra{font-size:.5rem;color:var(--cd);margin-top:.2rem;padding-top:.2rem;border-top:1px dashed var(--bg3);display:flex;gap:.4rem;flex-wrap:wrap}
.card-extra-pair{display:flex;gap:.15rem}
.card-extra-label{color:var(--cm);text-transform:uppercase;letter-spacing:.5px}
.card-extra-val{color:var(--cream)}

.empty-state{padding:3rem 2rem;text-align:center;color:var(--cm);font-style:italic;font-size:.8rem;width:100%}

.modal-bg{display:none;position:fixed;inset:0;background:rgba(0,0,0,.65);z-index:100;align-items:center;justify-content:center}
.modal-bg.open{display:flex}
.modal{background:var(--bg2);border:1px solid var(--bg3);padding:1.5rem;width:480px;max-width:92vw;max-height:90vh;overflow-y:auto}
.modal h2{font-size:.8rem;margin-bottom:1rem;color:var(--rust);letter-spacing:1px;text-transform:uppercase}
.fr{margin-bottom:.7rem}
.fr label{display:block;font-family:var(--mono);font-size:.55rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-bottom:.2rem}
.fr input,.fr select,.fr textarea{width:100%;padding:.4rem .5rem;background:var(--bg);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}
.fr input:focus,.fr select:focus,.fr textarea:focus{outline:none;border-color:var(--leather)}
.row2{display:grid;grid-template-columns:1fr 1fr;gap:.5rem}
.fr-section{margin-top:1rem;padding-top:.8rem;border-top:1px solid var(--bg3)}
.fr-section-label{font-family:var(--mono);font-size:.55rem;color:var(--rust);text-transform:uppercase;letter-spacing:1px;margin-bottom:.5rem}
.acts{display:flex;gap:.4rem;justify-content:flex-end;margin-top:1rem}
.acts .btn-del{margin-right:auto}
.stats-bar{padding:.4rem 1.2rem;border-top:1px solid var(--bg3);font-family:var(--mono);font-size:.5rem;color:var(--cm);display:flex;gap:.7rem;flex-wrap:wrap;background:var(--bg2);flex-shrink:0}
.stats-bar strong{color:var(--cd);font-weight:700}

.col-edit{padding:.4rem 0;display:flex;gap:.3rem;align-items:center}
.col-edit input{flex:1;padding:.3rem .4rem;background:var(--bg);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.65rem}
</style>
</head>
<body>

<div class="hdr">
<div class="hdr-left">
<h1 id="dash-title"><span>&#9670;</span> PRAIRIE</h1>
<select class="board-picker" id="board-picker" onchange="selectBoard(this.value)">
<option value="">Select board...</option>
</select>
</div>
<div class="hdr-actions">
<button class="btn" onclick="openBoardForm()">+ New Board</button>
<button class="btn" id="btn-edit-board" onclick="editCurrentBoard()" style="display:none">Edit Board</button>
</div>
</div>

<div class="board-area" id="board-area">
<div class="empty-state">Create your first board to get started.</div>
</div>

<div class="stats-bar" id="stats"></div>

<div class="modal-bg" id="mbg" onclick="if(event.target===this)closeModal()">
<div class="modal" id="mdl"></div>
</div>

<script>
var A='/api';
var boards=[],currentBoard=null,cards=[];
var boardExtras={},cardExtras={};
var cardCustomFields=[],boardCustomFields=[];
var dragCardID=null;

// ─── Loading ──────────────────────────────────────────────────────

async function loadAll(){
try{
var resps=await Promise.all([
fetch(A+'/boards').then(function(r){return r.json()}),
fetch(A+'/extras/boards').then(function(r){return r.json()}),
fetch(A+'/extras/cards').then(function(r){return r.json()}),
fetch(A+'/stats').then(function(r){return r.json()})
]);
boards=resps[0].boards||[];
boardExtras=resps[1]||{};
cardExtras=resps[2]||{};
renderStats(resps[3]||{});
renderBoardPicker();
}catch(e){
console.error('loadAll failed',e);
boards=[];
}
if(currentBoard){
var stillThere=null;
for(var i=0;i<boards.length;i++)if(boards[i].id===currentBoard.id){stillThere=boards[i];break}
if(stillThere){
currentBoard=stillThere;
await loadCards();
}else{
currentBoard=null;
renderBoard();
}
}else if(boards.length){
selectBoard(boards[0].id);
}else{
renderBoard();
}
}

async function loadCards(){
if(!currentBoard){cards=[];renderBoard();return}
try{
var r=await fetch(A+'/boards/'+currentBoard.id+'/cards').then(function(r){return r.json()});
cards=r.cards||[];
cards.forEach(function(c){
var x=cardExtras[c.id];
if(!x)return;
Object.keys(x).forEach(function(k){if(c[k]===undefined)c[k]=x[k]});
});
}catch(e){
cards=[];
}
renderBoard();
}

function renderBoardPicker(){
var sel=document.getElementById('board-picker');
var current=currentBoard?currentBoard.id:'';
var html='<option value="">Select board...</option>';
boards.forEach(function(b){
html+='<option value="'+esc(b.id)+'"'+(b.id===current?' selected':'')+'>'+esc(b.name)+' ('+b.card_count+')</option>';
});
sel.innerHTML=html;
}

function renderStats(s){
document.getElementById('stats').innerHTML=
'<span><strong>'+(s.boards||0)+'</strong> boards</span>'+
'<span><strong>'+(s.cards||0)+'</strong> cards</span>';
}

async function selectBoard(id){
if(!id){currentBoard=null;renderBoard();return}
for(var i=0;i<boards.length;i++)if(boards[i].id===id){currentBoard=boards[i];break}
document.getElementById('btn-edit-board').style.display=currentBoard?'inline-block':'none';
await loadCards();
}

function renderBoard(){
var area=document.getElementById('board-area');
if(!currentBoard){
area.innerHTML='<div class="empty-state">'+(boards.length?'Select a board from the dropdown above.':'Create your first board to get started.')+'</div>';
return;
}

var cols=currentBoard.columns||['Backlog','Todo','In Progress','Done'];
var byCol={};
cols.forEach(function(c){byCol[c]=[]});
cards.forEach(function(c){
if(byCol[c.column]===undefined)byCol[c.column]=[];
byCol[c.column].push(c);
});

var h='';
cols.forEach(function(col){
var colCards=byCol[col]||[];
h+='<div class="column" data-col="'+esc(col)+'" ondragover="onColDragOver(event,this)" ondragleave="onColDragLeave(this)" ondrop="onColDrop(event,this)">';
h+='<div class="col-hdr"><span class="col-name">'+esc(col)+'</span><span class="col-count">'+colCards.length+'</span></div>';
h+='<div class="col-cards">';
colCards.forEach(function(card){h+=cardHTML(card)});
h+='</div>';
h+='<div class="col-add"><input type="text" placeholder="+ Add card..." onkeydown="quickAddCard(event,\''+esc(col)+'\')"></div>';
h+='</div>';
});
area.innerHTML=h;
}

function cardHTML(c){
var h='<div class="card" draggable="true" data-id="'+esc(c.id)+'" ondragstart="onCardDragStart(event,\''+esc(c.id)+'\')" ondragend="onCardDragEnd(event)" onclick="openCardEdit(\''+esc(c.id)+'\')">';
h+='<div class="card-title">'+esc(c.title)+'</div>';

var hasMeta=c.assignee||c.labels;
if(hasMeta){
h+='<div class="card-meta">';
if(c.assignee)h+='<span class="card-assignee">'+esc(c.assignee)+'</span>';
if(c.labels){
String(c.labels).split(',').map(function(l){return l.trim()}).filter(function(l){return l}).forEach(function(l){
h+='<span class="card-label">'+esc(l)+'</span>';
});
}
h+='</div>';
}

// Custom fields
var customParts='';
cardCustomFields.forEach(function(f){
var v=c[f.name];
if(v===undefined||v===null||v==='')return;
customParts+='<span class="card-extra-pair"><span class="card-extra-label">'+esc(f.label)+':</span> <span class="card-extra-val">'+esc(String(v))+'</span></span>';
});
if(customParts)h+='<div class="card-extra">'+customParts+'</div>';

h+='</div>';
return h;
}

// ─── Drag and drop ────────────────────────────────────────────────

function onCardDragStart(ev,id){
dragCardID=id;
ev.dataTransfer.effectAllowed='move';
ev.target.classList.add('dragging');
}

function onCardDragEnd(ev){
ev.target.classList.remove('dragging');
dragCardID=null;
document.querySelectorAll('.column').forEach(function(c){c.classList.remove('drag-over')});
}

function onColDragOver(ev,col){
ev.preventDefault();
ev.dataTransfer.dropEffect='move';
col.classList.add('drag-over');
}

function onColDragLeave(col){
col.classList.remove('drag-over');
}

async function onColDrop(ev,col){
ev.preventDefault();
col.classList.remove('drag-over');
if(!dragCardID)return;
var newCol=col.getAttribute('data-col');
var card=null;
for(var i=0;i<cards.length;i++)if(cards[i].id===dragCardID){card=cards[i];break}
if(!card||card.column===newCol){return}
// Calculate position as last in target column
var maxPos=0;
cards.forEach(function(c){if(c.column===newCol&&c.position>maxPos)maxPos=c.position});
try{
await fetch(A+'/cards/'+dragCardID+'/move',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({column:newCol,position:maxPos+1})});
}catch(e){alert('Move failed');return}
await loadCards();
}

// ─── Quick add ────────────────────────────────────────────────────

async function quickAddCard(ev,col){
if(ev.key!=='Enter')return;
var inp=ev.target;
var title=inp.value.trim();
if(!title)return;
inp.disabled=true;
try{
await fetch(A+'/cards',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({board_id:currentBoard.id,title:title,column:col})});
inp.value='';
}catch(e){alert('Add failed')}
inp.disabled=false;
await loadCards();
inp.focus();
}

// ─── Card edit modal ──────────────────────────────────────────────

function openCardEdit(id){
var c=null;
for(var i=0;i<cards.length;i++)if(cards[i].id===id){c=cards[i];break}
if(!c)return;
var cols=currentBoard.columns||['Backlog','Todo','In Progress','Done'];

var h='<h2>Edit Card</h2>';
h+='<div class="fr"><label>Title *</label><input id="cf-title" value="'+esc(c.title)+'"></div>';
h+='<div class="fr"><label>Description</label><textarea id="cf-description" rows="3">'+esc(c.description||'')+'</textarea></div>';
h+='<div class="row2">';
h+='<div class="fr"><label>Column</label><select id="cf-column">';
cols.forEach(function(col){h+='<option value="'+esc(col)+'"'+(col===c.column?' selected':'')+'>'+esc(col)+'</option>'});
h+='</select></div>';
h+='<div class="fr"><label>Assignee</label><input id="cf-assignee" value="'+esc(c.assignee||'')+'"></div>';
h+='</div>';
h+='<div class="fr"><label>Labels (comma separated)</label><input id="cf-labels" value="'+esc(c.labels||'')+'"></div>';

if(cardCustomFields.length){
var label=window._cardCustomLabel||'Card Details';
h+='<div class="fr-section"><div class="fr-section-label">'+esc(label)+'</div>';
cardCustomFields.forEach(function(f){h+=customFieldHTML('xc',f,c[f.name])});
h+='</div>';
}

h+='<div class="acts">';
h+='<button class="btn btn-del" onclick="deleteCardConfirm(\''+esc(id)+'\')">Delete</button>';
h+='<button class="btn" onclick="closeModal()">Cancel</button>';
h+='<button class="btn btn-p" onclick="saveCard(\''+esc(id)+'\')">Save</button>';
h+='</div>';

document.getElementById('mdl').innerHTML=h;
document.getElementById('mbg').classList.add('open');
}

function customFieldHTML(prefix,f,value){
var v=value;
if(v===undefined||v===null)v='';
var h='<div class="fr"><label>'+esc(f.label)+'</label>';
if(f.type==='textarea'){
h+='<textarea id="'+prefix+'-'+f.name+'" rows="2">'+esc(String(v))+'</textarea>';
}else if(f.type==='select'){
h+='<select id="'+prefix+'-'+f.name+'"><option value="">Select...</option>';
(f.options||[]).forEach(function(o){
var sel=String(v)===String(o)?' selected':'';
h+='<option value="'+esc(String(o))+'"'+sel+'>'+esc(String(o))+'</option>';
});
h+='</select>';
}else if(f.type==='number'){
h+='<input type="number" id="'+prefix+'-'+f.name+'" value="'+esc(String(v))+'">';
}else{
h+='<input type="text" id="'+prefix+'-'+f.name+'" value="'+esc(String(v))+'">';
}
h+='</div>';
return h;
}

async function saveCard(id){
var title=document.getElementById('cf-title').value.trim();
if(!title){alert('Title required');return}
var body={
title:title,
description:document.getElementById('cf-description').value,
column:document.getElementById('cf-column').value,
assignee:document.getElementById('cf-assignee').value.trim(),
labels:document.getElementById('cf-labels').value.trim()
};
var extras={};
cardCustomFields.forEach(function(f){
var el=document.getElementById('xc-'+f.name);
if(!el)return;
extras[f.name]=f.type==='number'?(parseFloat(el.value)||0):el.value.trim();
});
try{
var r=await fetch(A+'/cards/'+id,{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});
if(!r.ok){var e=await r.json().catch(function(){return{}});alert(e.error||'Save failed');return}
if(Object.keys(extras).length){
await fetch(A+'/extras/cards/'+id,{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify(extras)}).catch(function(){});
}
}catch(e){alert('Network error');return}
closeModal();
await loadAll();
}

async function deleteCardConfirm(id){
if(!confirm('Delete this card?'))return;
await fetch(A+'/cards/'+id,{method:'DELETE'});
closeModal();
await loadAll();
}

// ─── Board form modal ─────────────────────────────────────────────

function openBoardForm(){
showBoardForm(null);
}

function editCurrentBoard(){
if(!currentBoard)return;
showBoardForm(currentBoard);
}

function showBoardForm(b){
var isEdit=!!b;
var board=b||{name:'',description:'',columns:['Backlog','Todo','In Progress','Done']};
var ext=isEdit?(boardExtras[board.id]||{}):{};

var h='<h2>'+(isEdit?'Edit Board':'New Board')+'</h2>';
h+='<div class="fr"><label>Name *</label><input id="bf-name" value="'+esc(board.name)+'"></div>';
h+='<div class="fr"><label>Description</label><textarea id="bf-description" rows="2">'+esc(board.description||'')+'</textarea></div>';
h+='<div class="fr"><label>Columns (comma separated)</label><input id="bf-columns" value="'+esc((board.columns||[]).join(', '))+'"></div>';

if(boardCustomFields.length){
var label=window._boardCustomLabel||'Board Details';
h+='<div class="fr-section"><div class="fr-section-label">'+esc(label)+'</div>';
boardCustomFields.forEach(function(f){h+=customFieldHTML('xb',f,ext[f.name])});
h+='</div>';
}

h+='<div class="acts">';
if(isEdit)h+='<button class="btn btn-del" onclick="deleteBoardConfirm(\''+esc(board.id)+'\')">Delete</button>';
h+='<button class="btn" onclick="closeModal()">Cancel</button>';
h+='<button class="btn btn-p" onclick="saveBoard(\''+(isEdit?esc(board.id):'')+'\')">'+(isEdit?'Save':'Create')+'</button>';
h+='</div>';

document.getElementById('mdl').innerHTML=h;
document.getElementById('mbg').classList.add('open');
setTimeout(function(){var n=document.getElementById('bf-name');if(n)n.focus()},50);
}

async function saveBoard(id){
var name=document.getElementById('bf-name').value.trim();
if(!name){alert('Name required');return}
var colsRaw=document.getElementById('bf-columns').value;
var cols=colsRaw.split(',').map(function(c){return c.trim()}).filter(function(c){return c});
if(!cols.length)cols=['Backlog','Todo','In Progress','Done'];
var body={
name:name,
description:document.getElementById('bf-description').value,
columns:cols
};
var extras={};
boardCustomFields.forEach(function(f){
var el=document.getElementById('xb-'+f.name);
if(!el)return;
extras[f.name]=f.type==='number'?(parseFloat(el.value)||0):el.value.trim();
});

var savedId=id;
try{
if(id){
var r1=await fetch(A+'/boards/'+id,{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});
if(!r1.ok){var e1=await r1.json().catch(function(){return{}});alert(e1.error||'Save failed');return}
}else{
var r2=await fetch(A+'/boards',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});
if(!r2.ok){var e2=await r2.json().catch(function(){return{}});alert(e2.error||'Create failed');return}
var created=await r2.json();
savedId=created.id;
}
if(savedId&&Object.keys(extras).length){
await fetch(A+'/extras/boards/'+savedId,{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify(extras)}).catch(function(){});
}
}catch(e){alert('Network error');return}
closeModal();
await loadAll();
if(!id&&savedId)selectBoard(savedId);
}

async function deleteBoardConfirm(id){
if(!confirm('Delete this board and all its cards?'))return;
await fetch(A+'/boards/'+id,{method:'DELETE'});
if(currentBoard&&currentBoard.id===id)currentBoard=null;
closeModal();
await loadAll();
}

function closeModal(){
document.getElementById('mbg').classList.remove('open');
}

function esc(s){
if(s===undefined||s===null)return'';
var d=document.createElement('div');
d.textContent=String(s);
return d.innerHTML;
}

document.addEventListener('keydown',function(e){if(e.key==='Escape')closeModal()});

// ─── Personalization ──────────────────────────────────────────────

(function loadPersonalization(){
fetch('/api/config').then(function(r){return r.json()}).then(function(cfg){
if(!cfg||typeof cfg!=='object')return;

if(cfg.dashboard_title){
var h1=document.getElementById('dash-title');
if(h1)h1.innerHTML='<span>&#9670;</span> '+esc(cfg.dashboard_title);
document.title=cfg.dashboard_title;
}

if(cfg.card_section_label)window._cardCustomLabel=cfg.card_section_label;
if(cfg.board_section_label)window._boardCustomLabel=cfg.board_section_label;

if(Array.isArray(cfg.card_custom_fields)){
cardCustomFields=cfg.card_custom_fields.filter(function(f){return f&&f.name&&f.label});
}
if(Array.isArray(cfg.board_custom_fields)){
boardCustomFields=cfg.board_custom_fields.filter(function(f){return f&&f.name&&f.label});
}
}).catch(function(){
}).finally(function(){
loadAll();
});
})();
</script>
</body>
</html>`
