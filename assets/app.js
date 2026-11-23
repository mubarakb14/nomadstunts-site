/* Nomad Stunts — shared header/footer, i18n (kk/ru/en) and page renderers. Content lives in data.js. */
(function(){
const D = window.NS_DATA;
const IMG = '.';
const UP = IMG + '/upload/images/';
const SR = D.showreel || {};                       // showreel: uploaded mp4, else YouTube id
const C = D.contacts || {phones:[]};               // phones, email, training base, map link
const LANGS = [['kk','ҚАЗ'],['ru','РУС'],['en','ENG']];

const I18N = window.NS_I18N;                      // texts: assets/i18n.js
const T = I18N.t, SERVICES = I18N.services, DIRECTIONS = I18N.directions, HISTORY = I18N.history, STARS = I18N.stars;


const ICONS = {
 ig:'<path d="M12 2.2c3.2 0 3.6 0 4.8.1 3.3.1 4.8 1.7 4.9 4.9.1 1.3.1 1.6.1 4.8s0 3.6-.1 4.8c-.1 3.2-1.7 4.8-4.9 4.9-1.3.1-1.6.1-4.8.1s-3.6 0-4.8-.1c-3.3-.1-4.8-1.7-4.9-4.9C2.2 15.6 2.2 15.2 2.2 12s0-3.6.1-4.8C2.4 3.9 3.9 2.4 7.2 2.3 8.4 2.2 8.8 2.2 12 2.2zM12 0C8.7 0 8.3 0 7.1.1 2.7.3.3 2.7.1 7.1 0 8.3 0 8.7 0 12s0 3.7.1 4.9c.2 4.4 2.6 6.8 7 7 1.2.1 1.6.1 4.9.1s3.7 0 4.9-.1c4.4-.2 6.8-2.6 7-7 .1-1.2.1-1.6.1-4.9s0-3.7-.1-4.9c-.2-4.4-2.6-6.8-7-7C15.7 0 15.3 0 12 0zm0 5.8a6.2 6.2 0 1 0 0 12.4 6.2 6.2 0 0 0 0-12.4zM12 16a4 4 0 1 1 0-8 4 4 0 0 1 0 8zm6.4-11.8a1.4 1.4 0 1 0 0 2.9 1.4 1.4 0 0 0 0-2.9z"/>',
 fb:'<path d="M24 12a12 12 0 1 0-13.9 11.9v-8.4H7.1V12h3V9.4c0-3 1.8-4.7 4.6-4.7 1.3 0 2.7.2 2.7.2v3h-1.5c-1.5 0-2 .9-2 1.9V12h3.4l-.5 3.5h-2.9v8.4A12 12 0 0 0 24 12z"/>',
 yt:'<path d="M23.5 6.2a3 3 0 0 0-2.1-2.1C19.5 3.6 12 3.6 12 3.6s-7.5 0-9.4.5A3 3 0 0 0 .5 6.2 31 31 0 0 0 0 12a31 31 0 0 0 .5 5.8 3 3 0 0 0 2.1 2.1c1.9.5 9.4.5 9.4.5s7.5 0 9.4-.5a3 3 0 0 0 2.1-2.1A31 31 0 0 0 24 12a31 31 0 0 0-.5-5.8zM9.6 15.6V8.4l6.3 3.6-6.3 3.6z"/>'
};
const SOCIAL = [['ig','Instagram','https://www.instagram.com/nomadstunts/'],['fb','Facebook','https://www.facebook.com/nomadstunts'],['yt','YouTube','https://www.youtube.com/channel/UCx2yGDE7QSX8n6cWHRJ2i7g']];

// ---------- language ----------
let lang = 'en';
try{
  const q = new URLSearchParams(location.search).get('lang');
  const s = q || localStorage.getItem('ns_lang');
  if (s && LANGS.some(l=>l[0]===s)) { lang = s; if (q) localStorage.setItem('ns_lang', q); }
  else { const n=(navigator.language||'').toLowerCase(); lang = n.startsWith('kk')?'kk':n.startsWith('ru')?'ru':'en'; }
}catch(e){}
const t = (k, v) => { const e=T[k]; let s = e ? (e[lang]||e.en) : k; if (v) for (const x in v) s=s.replace('{'+x+'}',v[x]); return s; };
const L = o => o ? (o[lang] || o.ru || o.en || '') : '';
const esc = s => String(s==null?'':s).replace(/[&<>"]/g,c=>({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;'}[c]));
const $ = (s,r=document) => r.querySelector(s);
const $$ = (s,r=document) => [...r.querySelectorAll(s)];
const img = p => p.startsWith('http') ? p : (p.startsWith('/') ? IMG + p : UP + p);

// ---------- chrome ----------
const PAGE = document.body.dataset.page;
const NAV = [['about','about'],['team','team'],['projects','projects'],['gallery','gallery'],['contacts','contacts']];  // clean URLs (server and Vercel map them to *.html)
const phones = cls => (C.phones||[]).map(p=>`<a${cls?` class="${cls}"`:''} href="tel:${p.replace(/[^\d+]/g,'')}">${esc(p)}</a>`).join('');
const mailto = () => C.email ? 'mailto:'+encodeURI(C.email) : '';
const soc = () => SOCIAL.map(([k,n,u])=>`<a href="${u}" target="_blank" rel="noopener" aria-label="${n}"><svg viewBox="0 0 24 24">${ICONS[k]}</svg></a>`).join('');

function chrome(){
  const hdr = document.createElement('header');
  hdr.innerHTML = `<a href="./" class="logo"><img src="assets/logo.svg" alt="Nomad Stunts" width="167" height="76"></a>
  <nav id="nav">${NAV.map(([k,u])=>`<a href="${u}" class="${k===PAGE?'on':''}">${t('nav_'+k)}</a>`).join('')}</nav>
  <div class="hdr-r"><div class="soc">${soc()}</div>
  <div class="lang" role="group" aria-label="Language">${LANGS.map(([c,n])=>`<button data-l="${c}" class="${c===lang?'on':''}" lang="${c}">${n}</button>`).join('')}</div></div>
  <button class="burger" aria-label="Menu"><i></i></button>`;
  const old = $('header'); old ? old.replaceWith(hdr) : document.body.prepend(hdr);

  const ft = document.createElement('footer'); ft.id='footer';
  ft.innerHTML = `<div class="wrap ft">
   <div><a href="./" class="ft-logo"><img src="assets/logo.svg" alt="Nomad Stunts" width="167" height="76"></a><h4>${t('ft_menu')}</h4><ul><li><a href="./">${t('home')}</a></li>${NAV.map(([k,u])=>`<li><a href="${u}">${t('nav_'+k)}</a></li>`).join('')}</ul></div>
   <div><h4>${t('ft_contact')}</h4><p>${t('ft_p')}</p>
    <div class="ph">${phones()}</div>
    ${C.email ? `<p style="margin-top:8px">Email: <a href="${mailto()}">${esc(C.email)}</a></p>` : ''}</div>
   <div class="ft-r"><div class="soc">${soc()}</div>${C.email ? `<a class="btn" href="${mailto()}">${t('send')}</a>` : ''}</div>
  </div><p class="copy">© 2000–${new Date().getFullYear()} Nomad Stunts · ${t('tagline')}. ${t('rights')}.</p>`;
  const oldf = $('#footer'); oldf ? oldf.replaceWith(ft) : document.body.append(ft);

  $$('.lang button',hdr).forEach(b=>b.onclick=()=>setLang(b.dataset.l));
  $('.burger',hdr).onclick=()=>document.body.classList.toggle('menu');
  document.documentElement.lang = lang;
}
function setLang(l){
  lang = l; try{localStorage.setItem('ns_lang',l)}catch(e){}
  const u = new URL(location.href); u.searchParams.delete('lang'); history.replaceState(null,'',u);
  render();
}
function i18nStatic(){ $$('[data-i]').forEach(el=>{ el.innerHTML = t(el.dataset.i); }); }

// ---------- reveal ----------
const io = 'IntersectionObserver' in window ? new IntersectionObserver(es=>es.forEach(e=>{if(e.isIntersecting){e.target.classList.add('in');io.unobserve(e.target)}}),{threshold:.08}) : null;
function watch(){ $$('.rv:not(.in)').forEach(el=> io ? io.observe(el) : el.classList.add('in')); }

// ---------- modal ----------
let modal, onKey;
function openModal(cls, html, keys){
  closeModal();
  modal = document.createElement('div');
  modal.className = 'modal open ' + cls;
  modal.setAttribute('role','dialog'); modal.setAttribute('aria-modal','true');
  modal.innerHTML = `<div class="box"><button class="x" aria-label="${t('close')}">×</button>${html}</div>`;
  modal.onclick = e => { if (e.target===modal || e.target.closest('.x')) closeModal(); };
  document.body.append(modal); document.body.classList.add('lock');
  onKey = e => { if (e.key==='Escape') closeModal(); else if (keys && keys[e.key]) keys[e.key](); };
  document.addEventListener('keydown', onKey);
  $('.x',modal).focus();
  return modal;
}
function closeModal(){
  if (!modal) return;
  modal.remove(); modal=null; document.body.classList.remove('lock');
  document.removeEventListener('keydown', onKey);
  if (location.hash) history.replaceState(null,'',location.pathname+location.search);
}
const hasReel = () => !!(SR.mp4 || SR.youtube);
const ytId = () => encodeURIComponent(SR.youtube||'');
function showreel(){
  const src = SR.mp4
    ? `<video src="${esc(img(SR.mp4))}" controls autoplay playsinline style="width:100%;height:100%"></video>`
    : `<iframe src="https://www.youtube-nocookie.com/embed/${ytId()}?autoplay=1&rel=0" allow="autoplay; encrypted-media; fullscreen" allowfullscreen title="Nomad Stunts Showreel"></iframe>`;
  openModal('vid', `<div class="vwrap">${src}</div>`);
}

// ---------- people & films ----------
const fullName = p => L(p.n).join(' ');
const roles = p => p.roles.map(r=>t('r_'+r)).join(', ');
function personModal(p, isFounder){
  const hon = L(p.h);
  const role = isFounder ? L(p.role) : roles(p);
  openModal('person', `<div class="mimg"><img src="${img(p.photo||p.img)}" alt="${esc(fullName(p))}"></div>
   <div class="mtxt">${isFounder ? `<div class="eyebrow">${t('cat_founder')}</div>` : ''}
   <h2 class="h">${esc(L(p.n)[0])}<br>${esc(L(p.n)[1])}</h2><p class="role">${esc(role)}</p>
   ${hon ? `<h4>${t('honours')}</h4><p>${esc(hon)}</p>` : ''}</div>`);
}
const memCard = (p,i) => `<button class="mem" data-p="${i}"><div class="im"><img src="${img(p.img)}" alt="${esc(fullName(p))}" loading="lazy"></div><div class="nm">${esc(fullName(p))}<small>${esc(roles(p))}</small></div></button>`;

const filmTitle = f => f.t[lang] || f.t.en || f.t.ru;
const filmInfo = f => { const o = f[lang] && Object.keys(f[lang]).length ? f[lang] : (f.en && Object.keys(f.en).length ? f.en : f.ru); return o||{}; };
const posterCard = f => `<button class="poster" data-f="${f.id}"><img src="${img(f.poster)}" alt="${esc(filmTitle(f))}" loading="lazy"><span class="cap">${esc(filmTitle(f))}<em>${f.y||''}</em></span></button>`;
function filmModal(id){
  const f = D.films.find(x=>x.id===id); if (!f) return;
  const o = filmInfo(f);
  const row = k => o[k] ? `<div><dt>${t(k)}</dt><dd>${esc(o[k])}</dd></div>` : '';
  openModal('film', `<div class="mimg"><img src="${img(f.img)}" alt="${esc(filmTitle(f))}"></div>
   <div class="mtxt"><div class="eyebrow">${o.type ? esc(o.type)+' · ' : ''}${f.y||''}</div><h2 class="h">${esc(filmTitle(f))}</h2>
   <dl class="facts">${row('country')}${row('director')}${row('stars')}</dl></div>`);
  history.replaceState(null,'','#'+id);
}
function bindPosters(root){ root.onclick = e => { const b=e.target.closest('[data-f]'); if (b) filmModal(b.dataset.f); }; }
function orderedFilms(){ const fe = D.films.filter(f=>f.featured).sort((a,b)=>a.featured-b.featured); return [...fe, ...D.films.filter(f=>!f.featured)]; }

function reel(el){
  el.innerHTML = `<button class="arr p" aria-label="‹">&#8592;</button><div class="track">${orderedFilms().map(posterCard).join('')}</div><button class="arr n" aria-label="›">&#8594;</button>`;
  const tr = $('.track',el);
  $('.arr.n',el).onclick=()=>tr.scrollBy({left:tr.clientWidth*.8,behavior:'smooth'});
  $('.arr.p',el).onclick=()=>tr.scrollBy({left:-tr.clientWidth*.8,behavior:'smooth'});
  bindPosters(tr);
}
function marquee(el){ const s = STARS[lang].map(x=>`${esc(x)}<i>&#9670;</i>`).join(''); el.innerHTML = `<div>${s+s}</div>`; el.setAttribute('aria-label',t('worked')); }
function serviceCards(el){ el.innerHTML = SERVICES.map(s=>`<article class="card rv"><h3>${s[lang][0]}</h3><div class="ph"><img src="${img(s.img)}" alt="" loading="lazy"></div><p>${s[lang][1]}</p></article>`).join(''); }
// FAQ answers can use live content: {projects}, {base}, {phones}, {email}
const fill = s => s.replace(/\{(projects|base|phones|email)\}/g, (_,k) => esc({projects:D.films.length, base:L(C.base), phones:(C.phones||[]).join(', '), email:C.email||''}[k]));
function faq(el){
  if (!el) return;
  el.innerHTML = ((I18N.faq||{})[lang]||[]).map(([q,a])=>`<details class="rv"><summary>${esc(q)}</summary><p>${fill(esc(a))}</p></details>`).join('');
}
function timeline(el){ el.innerHTML = HISTORY[lang].map(([y,x])=>`<li class="rv"><b>${y}</b><p>${x}</p></li>`).join(''); }
function stats(el){ el.innerHTML = `<div class="stat"><b>2000</b><span>${t('st1')}</span></div><div class="stat"><b>86</b><span>${t('st2')}</span></div><div class="stat"><b>40</b><span>${t('st3')}</span></div><div class="stat"><b>${D.films.length}+</b><span>${t('st4')}</span></div>`; }
function founderBlock(el){
  el.innerHTML = `<div class="img" role="img" aria-label="${esc(fullName(D.founder))}"></div>
  <div class="txt rv in"><div class="eyebrow">${t('f_eye')}</div><h2 class="h">${esc(fullName(D.founder))}</h2><p>${t('f_p')}</p>
  <div class="awards"><span>Taurus World Stunt Awards 2020</span><span>${t('aw2')}</span><span>${t('aw3')}</span><span>Hitz International</span></div>
  <button class="btn" data-founder>${t('profile')}</button></div>`;
  $('[data-founder]',el).onclick=()=>personModal(D.founder,true);
}

// ---------- pages ----------
let heroTimer;
const PAGES = {
 home(){
  const box = $('#slides');
  if (!box.dataset.ready){
    box.dataset.ready = 1;
    const HERO = ['big/44798_710250_00.webp','21109_783683_14.webp','big/97241_93672_00.webp'];
    HERO.forEach((h,i)=>{ const d=document.createElement('div'); d.className='sl'+(i?'':' on'); d.style.backgroundImage=`url(${img(h)})`; box.prepend(d);
      const b=document.createElement('button'); b.setAttribute('aria-label','Slide '+(i+1)); if(!i) b.className='on'; b.onclick=()=>go(i); $('#dots').append(b); });
    const slides = $$('.sl',box).reverse(); let cur=0;
    const go = i => { slides[cur].classList.remove('on'); $('#dots').children[cur].classList.remove('on'); cur=i; slides[cur].classList.add('on'); $('#dots').children[cur].classList.add('on'); clearInterval(heroTimer); heroTimer=setInterval(()=>go((cur+1)%slides.length),6000); };
    go(0);
    // muted background showreel over the slideshow
    if (!matchMedia('(prefers-reduced-motion: reduce)').matches && (SR.mp4 || (SR.youtube && location.protocol!=='file:'))){
      const v = document.createElement('div'); v.className='vid';
      v.innerHTML = SR.mp4
        ? `<video src="${esc(img(SR.mp4))}" autoplay muted loop playsinline style="width:100%;height:100%;object-fit:cover"></video>`
        : `<iframe src="https://www.youtube-nocookie.com/embed/${ytId()}?autoplay=1&mute=1&loop=1&playlist=${ytId()}&controls=0&rel=0&playsinline=1&iv_load_policy=3&disablekb=1" allow="autoplay; encrypted-media" tabindex="-1" aria-hidden="true" title="Showreel"></iframe>`;
      $('.dots',box).before(v);
      const show = () => setTimeout(()=>{ v.classList.add('on'); clearInterval(heroTimer); $('#dots').style.display='none'; }, SR.mp4?0:1800);
      (v.firstElementChild).addEventListener(SR.mp4?'playing':'load', show, {once:true});
    }
    if (hasReel()) $('#play').onclick = showreel; else $('#play').remove();
  }
  stats($('#stats')); serviceCards($('#cards')); founderBlock($('#founder'));
  $('#grid').innerHTML = D.team.slice(0,12).map(memCard).join('');
  $('#grid').onclick = e => { const b=e.target.closest('[data-p]'); if(b) personModal(D.team[b.dataset.p]); };
  reel($('#reel')); marquee($('#stars'));
 },
 about(){
  timeline($('#tl'));
  $('#dirs').innerHTML = DIRECTIONS.map(d=>`<div class="rv"><i class="ico" aria-hidden="true" style="-webkit-mask-image:url(${img(d[0])});mask-image:url(${img(d[0])})"></i><span>${esc({en:d[1],ru:d[2],kk:d[3]}[lang])}</span></div>`).join('');
  founderBlock($('#founder')); stats($('#stats')); marquee($('#stars')); faq($('#faq'));
  const v=$('#video');
  if (hasReel()) { if (SR.poster) v.style.backgroundImage=`url("${img(SR.poster)}")`; v.onclick=showreel; $('#video .lbl').textContent = t('watch'); }
  else v.closest('section').remove();
 },
 team(){
  const q = ($('#q').value||'').trim().toLowerCase();
  $('#tabs').innerHTML = `<button class="on">${t('all_team')}<sup>${D.team.length+1}</sup></button><button data-founder>${t('cat_founder')}</button>`;
  $('#tabs').onclick = e => { if (e.target.closest('[data-founder]')) personModal(D.founder,true); };
  $('#q').placeholder = t('search'); $('#q').oninput = PAGES.team;
  const match = p => !q || ['en','ru'].some(l=>p.n[l].join(' ').toLowerCase().includes(q));
  const list = D.team.filter(match);
  const showF = match(D.founder);
  const html = (showF ? `<button class="mem" data-founder><div class="im"><img src="${img(D.founder.photo)}" alt="${esc(fullName(D.founder))}"></div><div class="nm">${esc(fullName(D.founder))}<small>${t('cat_founder')}</small></div></button>` : '') +
    list.map(p=>memCard(p,D.team.indexOf(p))).join('');
  $('#grid').innerHTML = html || `<p class="empty">${t('empty')}</p>`;
  $('#grid').onclick = e => { if (e.target.closest('[data-founder]')) return personModal(D.founder,true); const b=e.target.closest('[data-p]'); if(b) personModal(D.team[b.dataset.p]); };
  $('#count').textContent = t('n_ppl',{n:list.length+(showF?1:0)});
 },
 projects(){
  const st = $('#years').dataset; const cur = st.y || 'all';
  const years = [...new Set(D.films.map(f=>f.y).filter(Boolean))].sort((a,b)=>b-a);
  $('#years').innerHTML = `<button data-y="all" class="${cur==='all'?'on':''}">${t('all')}</button>` + years.map(y=>`<button data-y="${y}" class="${String(y)===cur?'on':''}">${y}</button>`).join('');
  $('#years').onclick = e => { const b=e.target.closest('button'); if(b){ st.y=b.dataset.y; PAGES.projects(); } };
  const list = cur==='all' ? D.films : D.films.filter(f=>String(f.y)===cur);
  $('#count').textContent = t('n_proj',{n:list.length});
  $('#pgrid').innerHTML = list.map(posterCard).join('');
  bindPosters($('#pgrid'));
  if (location.hash && !modal) filmModal(decodeURIComponent(location.hash.slice(1)));
  marquee($('#stars'));
 },
 gallery(){
  const g = $('#gal');
  if (!g.children.length) g.innerHTML = D.gallery.map((x,i)=>`<button data-g="${i}" class="rv"><img src="${esc(img(x.thumb||x.img))}" data-full="${esc(img(x.img))}" alt="Nomad Stunts — ${esc(t('gal_p'))} (${i+1})" loading="lazy" onerror="this.onerror=null;this.src=this.dataset.full"></button>`).join('');
  g.onclick = e => { const b=e.target.closest('[data-g]'); if (b) lightbox(+b.dataset.g); };
 },
 contacts(){
  $('#phones').innerHTML = phones('big');
  $('#email').innerHTML = C.email ? `<a class="big" href="${mailto()}">${esc(C.email)}</a>` : '';
  $('#sendmail').href = mailto();
  $('#base').textContent = L(C.base);
  const m = $('#map'); if (C.map) m.href = C.map; else m.parentElement.remove();
 }
};
function lightbox(i){
  const n = D.gallery.length;
  const show = j => { i=(j+n)%n; $('img',modal).src = img(D.gallery[i].img); };
  openModal('lb', `<img alt=""><button class="arr p" aria-label="‹">&#8592;</button><button class="arr n" aria-label="›">&#8594;</button>`, {ArrowLeft:()=>show(i-1), ArrowRight:()=>show(i+1)});
  $('.arr.p',modal).onclick = e => { e.stopPropagation(); show(i-1); };
  $('.arr.n',modal).onclick = e => { e.stopPropagation(); show(i+1); };
  show(i);
}

function render(){
  chrome(); i18nStatic();
  const seo = (I18N.seo||{})[PAGE];
  document.title = seo ? (seo[lang]||seo.en)[0] : 'Nomad Stunts';
  (PAGES[PAGE]||(()=>{}))();
  watch();
}
document.addEventListener('click', e => { if (e.target.closest('nav a')) document.body.classList.remove('menu'); });
render();
})();
