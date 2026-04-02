package server
import "net/http"
func(s *Server)dashboard(w http.ResponseWriter,r *http.Request){w.Header().Set("Content-Type","text/html; charset=utf-8");w.Write([]byte(dashHTML))}
const dashHTML=`<!DOCTYPE html><html><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1"><title>Prairie</title>
<style>:root{--bg:#1a1410;--bg2:#241e18;--bg3:#2e261e;--rust:#c45d2c;--rl:#e8753a;--leather:#a0845c;--cream:#f0e6d3;--cd:#bfb5a3;--cm:#7a7060;--gold:#d4a843;--green:#4a9e5c;--blue:#4a7ec4;--mono:'JetBrains Mono',Consolas,monospace;--serif:'Libre Baskerville',Georgia,serif}*{margin:0;padding:0;box-sizing:border-box}body{background:var(--bg);color:var(--cream);font-family:var(--mono);font-size:13px;line-height:1.6;height:100vh;overflow:hidden}
.hdr{padding:.6rem 1.2rem;border-bottom:1px solid var(--bg3);display:flex;justify-content:space-between;align-items:center}.hdr h1{font-family:var(--serif);font-size:1rem}.hdr h1 span{color:var(--rl)}
.hdr-right{display:flex;gap:.5rem;align-items:center}
.btn{font-family:var(--mono);font-size:.68rem;padding:.3rem .6rem;border:1px solid;cursor:pointer;background:transparent}.btn-p{border-color:var(--rust);color:var(--rl)}.btn-p:hover{background:var(--rust);color:var(--cream)}
.board{display:flex;gap:.6rem;padding:1rem;overflow-x:auto;height:calc(100vh - 45px)}
.column{width:250px;flex-shrink:0;background:var(--bg2);border:1px solid var(--bg3);border-radius:4px;display:flex;flex-direction:column;max-height:100%}
.col-hdr{padding:.5rem .6rem;border-bottom:1px solid var(--bg3);font-size:.75rem;font-weight:600;display:flex;justify-content:space-between;align-items:center}
.col-count{font-size:.6rem;color:var(--cm);font-weight:400}
.col-cards{flex:1;overflow-y:auto;padding:.4rem}
.card{background:var(--bg);border:1px solid var(--bg3);padding:.5rem;margin-bottom:.3rem;cursor:pointer;border-radius:3px;transition:.1s}.card:hover{border-color:var(--rust)}
.card-title{font-size:.75rem;font-weight:600;margin-bottom:.15rem}
.card-meta{font-size:.6rem;color:var(--cm);display:flex;gap:.4rem;flex-wrap:wrap}
.card-label{font-size:.5rem;padding:.05rem .2rem;background:var(--bg3);color:var(--ll);border-radius:2px}
.col-add{padding:.4rem .6rem;border-top:1px solid var(--bg3);text-align:center;font-size:.7rem;color:var(--cm);cursor:pointer}.col-add:hover{color:var(--rl)}
.empty-state{display:flex;align-items:center;justify-content:center;height:calc(100vh - 45px);color:var(--cm);font-style:italic;font-family:var(--serif)}
.modal-bg{position:fixed;top:0;left:0;right:0;bottom:0;background:rgba(0,0,0,.65);display:flex;align-items:center;justify-content:center;z-index:100}
.modal{background:var(--bg2);border:1px solid var(--bg3);padding:1.5rem;width:90%;max-width:450px}
.modal h2{font-family:var(--serif);font-size:.9rem;margin-bottom:1rem}
label.fl{display:block;font-size:.65rem;color:var(--leather);text-transform:uppercase;letter-spacing:1px;margin-bottom:.2rem;margin-top:.5rem}
input[type=text],textarea,select{background:var(--bg);border:1px solid var(--bg3);color:var(--cream);padding:.35rem .5rem;font-family:var(--mono);font-size:.78rem;width:100%;outline:none}textarea{resize:vertical;min-height:60px}
</style>
<link href="https://fonts.googleapis.com/css2?family=Libre+Baskerville:ital@0;1&family=JetBrains+Mono:wght@400;600&display=swap" rel="stylesheet">
</head><body>
<div class="hdr"><h1><span>Prairie</span></h1><div class="hdr-right"><select id="boardSelect" onchange="switchBoard(this.value)" style="background:var(--bg);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.72rem;padding:.25rem .4rem"></select><button class="btn btn-p" onclick="showNewBoard()">+ Board</button></div></div>
<div id="content"><div class="empty-state">Select or create a board</div></div>
<div id="modal"></div>
<script>
let boards=[],curBoard='',cards=[];
async function api(u,o){return(await fetch(u,o)).json()}
function esc(s){return String(s||'').replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;').replace(/"/g,'&quot;')}
async function init(){const d=await api('/api/boards');boards=d.boards||[];
document.getElementById('boardSelect').innerHTML='<option value="">Select board</option>'+boards.map(b=>'<option value="'+b.id+'">'+esc(b.name)+' ('+b.card_count+')</option>').join('');
if(curBoard)document.getElementById('boardSelect').value=curBoard}
function switchBoard(id){curBoard=id;loadBoard()}
async function loadBoard(){
if(!curBoard){document.getElementById('content').innerHTML='<div class="empty-state">Select or create a board</div>';return}
const[b,cd]=await Promise.all([api('/api/boards/'+curBoard),api('/api/boards/'+curBoard+'/cards')]);
cards=cd.cards||[];const cols=b.columns||[];
document.getElementById('content').innerHTML='<div class="board">'+cols.map(col=>{
const colCards=cards.filter(c=>c.column===col);
return '<div class="column"><div class="col-hdr">'+esc(col)+' <span class="col-count">'+colCards.length+'</span></div><div class="col-cards">'+
colCards.map(c=>{
const labels=(c.labels||'').split(',').filter(Boolean).map(l=>'<span class="card-label">'+esc(l.trim())+'</span>').join('');
return '<div class="card" onclick="showCard(\''+c.id+'\')"><div class="card-title">'+esc(c.title)+'</div><div class="card-meta">'+(c.assignee?'<span>'+esc(c.assignee)+'</span>':'')+labels+'</div></div>'}).join('')+
'</div><div class="col-add" onclick="addCard(\''+esc(col)+'\')">+ Add card</div></div>'}).join('')+'</div>'}
function addCard(col){
document.getElementById('modal').innerHTML='<div class="modal-bg" onclick="if(event.target===this)closeModal()"><div class="modal"><h2>New Card</h2><label class="fl">Title</label><input type="text" id="nc-title"><label class="fl">Description</label><textarea id="nc-desc" rows="2"></textarea><label class="fl">Assignee</label><input type="text" id="nc-assign"><label class="fl">Labels (comma-separated)</label><input type="text" id="nc-labels"><div style="display:flex;gap:.5rem;margin-top:1rem"><button class="btn btn-p" onclick="saveCard(\''+esc(col)+'\')">Create</button><button class="btn" style="border-color:var(--bg3);color:var(--cm)" onclick="closeModal()">Cancel</button></div></div></div>'}
async function saveCard(col){const b={board_id:curBoard,column:col,title:document.getElementById('nc-title').value,description:document.getElementById('nc-desc').value,assignee:document.getElementById('nc-assign').value,labels:document.getElementById('nc-labels').value};if(!b.title){alert('Title required');return};await api('/api/cards',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(b)});closeModal();loadBoard();init()}
async function showCard(id){const c=await api('/api/cards/'+id+'?');if(!c||c.error)return;
const board=boards.find(b=>b.id===curBoard);const cols=(board?board.columns:[])||[];
const colOpts=cols.map(col=>'<option'+(col===c.column?' selected':'')+'>'+esc(col)+'</option>').join('');
document.getElementById('modal').innerHTML='<div class="modal-bg" onclick="if(event.target===this)closeModal()"><div class="modal"><h2>'+esc(c.title)+'</h2>'+(c.description?'<div style="padding:.4rem;background:var(--bg);border:1px solid var(--bg3);font-size:.78rem;color:var(--cd);margin-bottom:.5rem">'+esc(c.description)+'</div>':'')+
'<label class="fl">Move to column</label><select id="mc-col">'+colOpts+'</select>'+
'<div style="display:flex;gap:.3rem;margin-top:1rem"><button class="btn btn-p" onclick="move(\''+c.id+'\')">Move</button><button class="btn" style="border-color:var(--bg3);color:var(--cm)" onclick="if(confirm(\'Delete?\'))delCard(\''+c.id+'\')">Del</button><button class="btn" style="border-color:var(--bg3);color:var(--cm)" onclick="closeModal()">Close</button></div></div></div>'}
async function move(id){const col=document.getElementById('mc-col').value;await api('/api/cards/'+id+'/move',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({column:col,position:0})});closeModal();loadBoard();init()}
async function delCard(id){await api('/api/cards/'+id,{method:'DELETE'});closeModal();loadBoard();init()}
function showNewBoard(){document.getElementById('modal').innerHTML='<div class="modal-bg" onclick="if(event.target===this)closeModal()"><div class="modal"><h2>New Board</h2><label class="fl">Name</label><input type="text" id="nb-name"><label class="fl">Columns (comma-separated)</label><input type="text" id="nb-cols" value="Backlog, Todo, In Progress, Done"><div style="display:flex;gap:.5rem;margin-top:1rem"><button class="btn btn-p" onclick="saveBoard()">Create</button><button class="btn" style="border-color:var(--bg3);color:var(--cm)" onclick="closeModal()">Cancel</button></div></div></div>'}
async function saveBoard(){const cols=(document.getElementById('nb-cols').value||'').split(',').map(s=>s.trim()).filter(Boolean);const b={name:document.getElementById('nb-name').value,columns:cols.length?cols:undefined};if(!b.name){alert('Name required');return};const r=await api('/api/boards',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(b)});curBoard=r.id;closeModal();init();loadBoard()}
function closeModal(){document.getElementById('modal').innerHTML=''}
init()
</script></body></html>`
