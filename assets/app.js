/* Nomad Stunts — shared header/footer, i18n (kk/ru/en) and page renderers. Content lives in data.js. */
(function(){
const D = window.NS_DATA;
const IMG = '.';
const UP = IMG + '/upload/images/';
// Showreel: YouTube id, or set SHOWREEL_MP4 to a video file URL to use it instead
const SHOWREEL = 'p89l3BQwIl8';
const SHOWREEL_MP4 = 'assets/video/showreel.mp4';
const LANGS = [['kk','ҚАЗ'],['ru','РУС'],['en','ENG']];

const T = {
nav_about:{kk:'Біз туралы',ru:'О нас',en:'About us'},
nav_team:{kk:'Команда',ru:'Команда',en:'Our team'},
nav_projects:{kk:'Жобалар',ru:'Проекты',en:'Projects'},
nav_gallery:{kk:'Галерея',ru:'Фотогалерея',en:'Gallery'},
nav_contacts:{kk:'Байланыс',ru:'Контакты',en:'Contacts'},
home:{kk:'Басты бет',ru:'Главная',en:'Home'},
tagline:{kk:'Халықаралық кәсіби каскадерлер тобы',ru:'Международная группа профессиональных каскадеров',en:'International Action Company'},
hero_p:{kk:'Халықаралық кәсіби каскадерлер тобы. 2000 жылдан бері',ru:'Международная группа профессиональных каскадеров. С 2000 года',en:'International Action Company. Est. 2000'},
watch:{kk:'Шоурилді көру',ru:'Смотреть шоурил',en:'Watch showreel'},
intro_t:{kk:'Экшн шеберлері',ru:'Мастера экшена',en:'Action specialists'},
intro_p:{
 kk:'Nomad Stunts — <strong>халықаралық кәсіби каскадерлер тобы</strong>. Оны 2000 жылы Taurus World Stunt Awards лауреаты, трюк қоюшы, каскадер, Қазақстан Республикасының еңбек сіңірген қайраткері <strong>Жайдарбек Кунгужинов</strong> құрды. Командада әлемнің түкпір-түкпірінен келген трюк қоюшылар, каскадерлер, хорс-мастерлер, риггерлер мен шайқас сахналарын қоюшылар бар.',
 ru:'Nomad Stunts — <strong>международная группа профессиональных каскадеров</strong>, которую в 2000 году организовал лауреат Taurus World Stunt Awards, постановщик трюков, каскадёр и заслуженный деятель Республики Казахстан <strong>Жайдарбек Кунгужинов</strong>. В команде — постановщики трюков, каскадёры, хорс-мастеры, риггеры и постановщики боёв со всего мира.',
 en:'Nomad Stunts is an <strong>international action company</strong> established in 2000 by the Taurus World Stunt Awards winner, stunt coordinator, stunt performer and Honored Artist of the Republic of Kazakhstan, <strong>Zhaidarbek Kunguzhinov</strong>. Our team unites stunt coordinators, stunt performers, horse masters, riggers and fight choreographers from around the world.'},
st1:{kk:'Құрылған жылы',ru:'Год основания',en:'Founded'},
st2:{kk:'Каскадер',ru:'Каскадёров',en:'Stunt performers'},
st3:{kk:'Үйретілген трюк жылқысы',ru:'Обученных трюковых лошадей',en:'Trained stunt horses'},
st4:{kk:'Кино және ТВ жобасы',ru:'Проектов в кино и на ТВ',en:'Film & TV projects'},
srv_t:{kk:'Біздің қызметтер',ru:'Наши услуги',en:'Our services'},
dir_t:{kk:'Жұмыс бағыттары',ru:'Направления работы',en:'What we do'},
hist_t:{kk:'Біздің тарих',ru:'Наша история',en:'Our history'},
hist_cap:{kk:'Жаттығу базасы, Алматы облысы, Есік қаласының маңы',ru:'Тренировочная база возле Есика, Алматинская область',en:'Training base near Issyk, Almaty Region'},
f_eye:{kk:'Негізін қалаушы',ru:'Основатель',en:'Founder'},
f_p:{kk:'Трюк қоюшы, каскадер, Қазақстан Республикасының еңбек сіңірген қайраткері. 2020 жылы «Томирис» фильміндегі экшн үшін каскадерлердің «Оскары» — Taurus World Stunt Awards сыйлығын алған алғашқы қазақстандық атанды.',
 ru:'Постановщик трюков, каскадёр, заслуженный деятель Республики Казахстан. В 2020 году стал первым казахстанцем, получившим Taurus World Stunt Awards — «Оскар» каскадёров — за экшен в фильме «Томирис».',
 en:'Stunt coordinator, stunt performer and Honored Artist of the Republic of Kazakhstan. In 2020 he became the first Kazakhstan national to win the Taurus World Stunt Awards — the stunt performers’ Oscars — for the action in “Tomiris”.'},
aw2:{kk:'Taurus номинанты 2017 · «Викинг»',ru:'Номинант Taurus 2017 · «Викинг»',en:'Taurus nominee 2017 · Viking'},
aw3:{kk:'«Құрмет» ордені',ru:'Орден «Құрмет»',en:'Kurmet Order'},
team_t:{kk:'Біздің команда',ru:'Наша команда',en:'Our team'},
team_all:{kk:'Бүкіл команда',ru:'Вся команда',en:'Meet the whole team'},
proj_t:{kk:'Жобалар',ru:'Проекты',en:'Projects'},
proj_all:{kk:'Барлық жобалар',ru:'Все проекты',en:'All projects'},
more_about:{kk:'Толығырақ',ru:'Подробнее о нас',en:'More about us'},
profile:{kk:'Толық профиль',ru:'Полный профиль',en:'Full profile'},
worked:{kk:'Біз жұмыс істеген жұлдыздар',ru:'Звёзды, с которыми мы работали',en:'Stars we have worked with'},
all:{kk:'Барлығы',ru:'Все',en:'All'},
n_proj:{kk:'Жобалар: {n}',ru:'Проектов: {n}',en:'{n} projects'},
n_ppl:{kk:'Табылды: {n}',ru:'Найдено: {n}',en:'{n} people'},
search:{kk:'Атымен іздеу',ru:'Поиск по имени',en:'Search by name'},
empty:{kk:'Ештеңе табылмады',ru:'Ничего не найдено',en:'Nothing found'},
skills:{kk:'Дағдылар',ru:'Навыки',en:'Skills'},
honours:{kk:'Марапаттар мен жетістіктер',ru:'Награды и достижения',en:'Honours and affiliations'},
year:{kk:'Жыл',ru:'Год',en:'Year'},
country:{kk:'Ел',ru:'Страна',en:'Country'},
director:{kk:'Режиссер',ru:'Режиссёр',en:'Director'},
stars:{kk:'Рөлдерде',ru:'В ролях',en:'Starring'},
cat_founder:{kk:'Негізін қалаушы',ru:'Основатель',en:'Founder'},
cat_coord:{kk:'Трюк қоюшылар',ru:'Постановщики трюков',en:'Stunt coordinators'},
cat_horse:{kk:'Хорс-мастерлер',ru:'Хорс-мастеры',en:'Horse masters'},
cat_rig:{kk:'Риггерлер',ru:'Риггеры',en:'Riggers'},
cat_perf:{kk:'Каскадерлер',ru:'Каскадёры',en:'Stunt performers'},
r_perf:{kk:'Каскадер',ru:'Каскадёр',en:'Stunt performer'},
r_coord:{kk:'Трюк қоюшы',ru:'Постановщик трюков',en:'Stunt coordinator'},
r_horse:{kk:'Хорс-мастер',ru:'Хорс-мастер',en:'Horse master'},
r_rig:{kk:'Риггер',ru:'Риггер',en:'Rigger'},
r_actor:{kk:'Актер',ru:'Актёр',en:'Actor'},
r_actress:{kk:'Актриса',ru:'Актриса',en:'Actress'},
r_cam:{kk:'Экшн-оператор',ru:'Экшен-оператор',en:'Action cameraman'},
gal_t:{kk:'Фотогалерея',ru:'Фотогалерея',en:'Gallery'},
gal_p:{kk:'Түсірілім алаңынан және жаттығулардан кадрлар',ru:'Кадры со съёмочных площадок и тренировок',en:'Moments from film sets and training'},
c_agent:{kk:'Nomad Stunts агенті',ru:'Агент Nomad Stunts',en:'Nomad Stunts agent'},
agent:{kk:'Мурат Елькенов',ru:'Мурат Елькенов',en:'Murat Yelkenov'},
c_write:{kk:'Бізге жазыңыз',ru:'Пишите нам',en:'Write to us'},
c_base:{kk:'Жаттығу базасы',ru:'Тренировочная база',en:'Training base'},
c_base_p:{kk:'Алматы облысы, Есік қаласының маңы. 40 үйретілген трюк жылқысы.',ru:'Алматинская область, недалеко от города Есик. 40 обученных трюковых лошадей.',en:'Near the town of Issyk, Almaty Region, Kazakhstan. Home to our 40 trained stunt horses.'},
c_soc:{kk:'Әлеуметтік желілер',ru:'Мы в соцсетях',en:'Follow us'},
c_lead:{kk:'Ынтымақтастық, түсірілім және кастинг сұрақтары бойынша агентімізге хабарласыңыз — біз әлемнің кез келген нүктесіндегі жобаларға шығамыз.',ru:'По вопросам сотрудничества, съёмок и кастинга обращайтесь к нашему агенту — мы работаем на проектах по всему миру.',en:'For collaborations, productions and casting, contact our agent — we work on projects all over the world.'},
ft_menu:{kk:'Мәзір',ru:'Меню',en:'Menu'},
ft_contact:{kk:'Бізбен байланыс',ru:'Свяжитесь с нами',en:'Contact us'},
ft_p:{kk:'Nomad Stunts-пен ынтымақтастық бойынша барлық сұрақтарға агентіміз жауап береді:',ru:'По всем вопросам сотрудничества с Nomad Stunts обращайтесь к нашему агенту:',en:'For all inquiries regarding Nomad Stunts please contact our agent:'},
send:{kk:'Хат жазу',ru:'Написать письмо',en:'Send email'},
rights:{kk:'Барлық құқықтар қорғалған',ru:'Все права защищены',en:'All rights reserved'},
close:{kk:'Жабу',ru:'Закрыть',en:'Close'}
};

const SERVICES = [
 {img:'big/70704_512807_01.jpg',
  en:['Stunt coordinators','Creative, safety-minded leaders who design large-scale action — from cavalry battles to precise close-quarter fights — and deliver it for international productions.'],
  ru:['Постановщики трюков','Креативные и ответственные лидеры, которые ставят масштабный экшен — от кавалерийских сражений до точных ближних боёв — для международных проектов.'],
  kk:['Трюк қоюшылар','Халықаралық жобалар үшін ауқымды экшнді — атты әскер шайқастарынан нақты жекпе-жек сахналарына дейін — шығармашылықпен әрі қауіпсіздікті ескеріп қоятын көшбасшылар.']},
 {img:'34364_20156_14.jpg',
  en:['Equestrian stunts','Our horse masters and 40 trained stunt horses deliver gallops, falls, horse acrobatics and battle sequences with stunning realism and full care for the animals.'],
  ru:['Конные трюки','Наши хорс-мастеры и 40 обученных трюковых лошадей — галоп, падения, джигитовка и батальные сцены с максимальной достоверностью и заботой о животных.'],
  kk:['Атпен орындалатын трюктер','Хорс-мастерлеріміз бен 40 үйретілген трюк жылқысы — шабыс, аттан құлау, жігіттік және шайқас сахналары барынша шынайы әрі жануарларға қамқорлықпен.']},
 {img:'82219_212255_14.jpg',
  en:['Stunt rigging','Riggers who engineer wire work, ratchets and descenders, making the impossible look effortless while keeping every performer secure.'],
  ru:['Риггинг','Риггеры, которые проектируют тросовые системы, рэтчеты и спуски — невозможное выглядит легко, а каждый каскадёр в безопасности.'],
  kk:['Риггинг','Арқан жүйелерін, рэтчеттер мен түсу құрылғыларын жобалайтын риггерлер — мүмкін емес нәрсе оңай көрінеді, ал әр каскадер қауіпсіз.']},
 {img:'36331_896275_14.jpg',
  en:['Fight choreography','Battle scenes with swords, spears, bows and bare hands — choreography that turns combat into visual storytelling.'],
  ru:['Постановка боёв','Батальные сцены с мечами, копьями, луками и без оружия — хореография, которая превращает бой в историю на экране.'],
  kk:['Шайқас сахналарын қою','Қылыш, найза, садақпен және қарусыз шайқас сахналары — ұрысты экрандағы әңгімеге айналдыратын хореография.']},
 {img:'68101_862130_14.jpg',
  en:['High falls & fireburns','High falls, including the fall of a burning man, full-body fireburns and pyrotechnics — executed under strict safety protocols.'],
  ru:['Падения и огонь','Падения с высоты, включая падение горящего человека, полные огненные трюки и пиротехника — по строгим протоколам безопасности.'],
  kk:['Құлау және от','Биіктен құлау, соның ішінде жанып жатқан адамның құлауы, толық жану трюктері және пиротехника — қатаң қауіпсіздік хаттамаларымен.']},
 {img:'big/97241_93672_00.jpg',
  en:['Stunt performers','Acrobatics, martial arts, archery, fencing and scenic movement — performers from Kazakhstan, Kyrgyzstan, Australia, Canada, Turkey and beyond.'],
  ru:['Каскадёры','Акробатика, единоборства, стрельба из лука, фехтование и сценическое движение — каскадёры из Казахстана, Кыргызстана, Австралии, Канады, Турции и других стран.'],
  kk:['Каскадерлер','Акробатика, жекпе-жек, садақ ату, семсерлесу және сахналық қозғалыс — Қазақстан, Қырғызстан, Австралия, Канада, Түркия және басқа елдерден келген каскадерлер.']}
];

const DIRECTIONS = [
 ['92148_134821_17.png','Equestrian stunts','Конные трюки','Атпен орындалатын трюктер'],
 ['23839_325489_17.png','Stunt rigging','Риггинг (работа с подвесными системами)','Риггинг (аспалы жүйелермен жұмыс)'],
 ['87155_850522_17.png','Acrobatics','Акробатика','Акробатика'],
 ['93344_329872_17.png','Fireburns','Горение','Жану трюктері'],
 ['28442_362218_17.png','Pyrotechnics, SFX','Пиротехника, SFX','Пиротехника, SFX'],
 ['66369_390104_17.png','Scenic movement','Сценическое движение','Сахналық қозғалыс'],
 ['53033_393953_17.png','Choreography and performance of battle scenes (with various types of weapons)','Постановка и исполнение боевых сцен (с разными видами оружия)','Шайқас сахналарын қою және орындау (әртүрлі қарумен)'],
 ['2239_714086_17.png','High falls (including fall of a burning man)','Падения с высоты (в том числе горящего человека)','Биіктен құлау (соның ішінде жанып жатқан адамның)']
];

const HITZ = '<a href="https://www.hitzinternational.com/" target="_blank" rel="noopener">Hitz International</a>';
const HISTORY = {
en:[
 ['2000','Nomad Stunts is an international action company, established in 2000 by the Honored Artist of the Republic of Kazakhstan, stunt coordinator and stunt performer Zhaidarbek Kunguzhinov.'],
 ['2003','At the invitation of Kazakhfilm studio, Nomad Stunts took part in the large-scale historical film “Nomad”. Zhaidarbek Kunguzhinov was invited to be assistant stunt coordinator and stunt double for Mark Dacascos and Jason Scott Lee.'],
 ['2007','“Mongol”, directed by Sergei Bodrov, was released, with Zhaidarbek Kunguzhinov as stunt coordinator, stunt performer and stunt double of the leading characters. Its action scenes received rave reviews from film industry professionals, and from that moment the Kazakhstan stunt team gained worldwide acknowledgement.'],
 ['2009','Since 2009 Nomad Stunts have been regular participants in large-scale national and Hollywood productions: Mulan, Tomiris, The Way Back, The Expendables 2, 47 Ronin, Conan the Barbarian, and the popular TV series Warrior, Marco Polo and Resurrection: Ertuğrul. The full filmography is <a href="projects.html">here</a>.'],
 ['2015','Zhaidarbek Kunguzhinov joined the '+HITZ+' Action Team, bringing Nomad Stunts into one of the world’s leading stunt teams.'],
 ['2017','Zhaidarbek Kunguzhinov was nominated for the Taurus World Stunt Awards for his work on “Viking”.'],
 ['2018','Zhaidarbek Kunguzhinov and Nomad Stunts were included in the KINES Kazakh Book of Records. Zhaidarbek also became a holder of the Kurmet Order, a state award of Kazakhstan.'],
 ['2020','Zhaidarbek Kunguzhinov made history by becoming the first Kazakhstan national to win the prestigious Taurus World Stunt Awards (USA), known as the stunt performers’ Oscars, for the action in “Tomiris”.'],
 ['Today','Over the years Nomad Stunts has worked with global stars such as Keanu Reeves, Colin Farrell, Joaquin Phoenix, Michelle Yeoh, Saoirse Ronan, Samuel L. Jackson, Natalie Portman, Liu Yifei, Sean Bean, Ed Harris, Bruce Willis, Jason Statham, Tadanobu Asano, Arnold Schwarzenegger, Jet Li, Donnie Yen, Tom Wu and Benedict Wong. Today the team is 86 stunt performers from Kazakhstan, Kyrgyzstan, Australia, Canada, Turkey, Azerbaijan, Russia and Ukraine, 40 trained stunt horses and a special training base near the town of Issyk in Almaty Region.']],
ru:[
 ['2000','Nomad Stunts — международная группа профессиональных каскадеров, которую в 2000 году организовал постановщик трюков, заслуженный деятель Республики Казахстан Жайдарбек Кунгужинов.'],
 ['2003','По приглашению студии «Казахфильм» каскадёры Nomad Stunts приняли участие в съёмках масштабной исторической картины «Кочевник». Жайдарбек Кунгужинов получил предложение стать ассистентом постановщика трюков, а также дублёром Марка Дакаскоса и Джейсона Скотта Ли.'],
 ['2007','На экраны вышел фильм Сергея Бодрова-старшего «Монгол», в котором Жайдарбек Кунгужинов выступил постановщиком трюков, каскадёром и дублёром главных героев. Зрелищные экшен-сцены получили высокую оценку профессионалов киноиндустрии, и с этого периода казахстанские каскадёры стали известны во всём мире.'],
 ['2009','С 2009 года каскадёры Nomad Stunts — постоянные участники создания масштабных отечественных и голливудских картин: «Мулан», «Томирис», «Путь домой», «Неудержимые 2», «47 ронинов», «Конан-варвар», «Мын бала», «Тобол». Полную фильмографию можно посмотреть <a href="projects.html">здесь</a>.'],
 ['2015','Жайдарбек Кунгужинов принят в Международную группу постановщиков трюков '+HITZ+'. Благодаря этому каскадёры Nomad Stunts вошли в состав одной из лучших мировых каскадёрских команд.'],
 ['2017','Жайдарбек Кунгужинов номинирован на премию Taurus World Stunt Awards за работу в фильме «Викинг».'],
 ['2018','Жайдарбек Кунгужинов и группа Nomad Stunts вошли в казахстанскую книгу рекордов КИНЭС. Жайдарбек Кунгужинов также стал обладателем государственной награды Республики Казахстан — ордена «Құрмет».'],
 ['2020','Жайдарбек Кунгужинов сотворил историю, став первым казахстанцем, получившим премию «Таурус» Мировой академии каскадёров (США). Эту престижнейшую награду называют каскадёрским «Оскаром». Премия вручена за постановку трюков в фильме «Томирис».'],
 ['Сегодня','За годы существования группы каскадёрам Nomad Stunts удалось поработать с такими звёздами мирового кино, как Киану Ривз, Колин Фаррелл, Хоакин Феникс, Мишель Йео, Сирша Ронан, Сэмюэл Л. Джексон, Натали Портман, Лю Ифэй, Шон Бин, Эд Харрис, Брюс Уиллис, Джейсон Стейтем, Таданобу Асано, Арнольд Шварценеггер, Джет Ли, Донни Йен, Том Ву, Бенедикт Вонг. Сегодня Nomad Stunts — это 86 каскадёров из Казахстана, Кыргызстана, Австралии, Канады, Турции, Азербайджана, России и Украины, 40 обученных трюковых лошадей и специальная тренировочная база недалеко от города Есик в Алматинской области.']],
kk:[
 ['2000','Nomad Stunts — халықаралық кәсіби каскадерлер тобы. Оны 2000 жылы трюк қоюшы, Қазақстан Республикасының еңбек сіңірген қайраткері Жайдарбек Кунгужинов құрды.'],
 ['2003','«Қазақфильм» киностудиясының шақыруымен Nomad Stunts каскадерлері «Көшпенділер» атты ауқымды тарихи фильмнің түсірілініне қатысты. Жайдарбек Кунгужиновке трюк қоюшының ассистенті, сондай-ақ Марк Дакаскос пен Джейсон Скотт Лидің дублері болу ұсынылды.'],
 ['2007','Сергей Бодров-үлкеннің «Моңғол» фильмі экранға шықты: онда Жайдарбек Кунгужинов трюк қоюшы, каскадер және басты кейіпкерлердің дублері болды. Әсерлі экшн-сахналар кино индустриясы мамандарының жоғары бағасын алып, осы кезеңнен бастап қазақстандық каскадерлер бүкіл әлемге танылды.'],
 ['2009','2009 жылдан бері Nomad Stunts каскадерлері ауқымды отандық және Голливуд фильмдерінің тұрақты қатысушылары: «Мулан», «Томирис», «Үйге жол», «Қайсарлар 2», «47 ронин», «Конан-варвар», «Мың бала», «Тобыл». Толық фильмографияны <a href="projects.html">мұнда</a> көруге болады.'],
 ['2015','Жайдарбек Кунгужинов '+HITZ+' халықаралық трюк қоюшылар тобына қабылданды. Осылайша Nomad Stunts каскадерлері әлемдегі үздік каскадерлер командаларының біріне енді.'],
 ['2017','Жайдарбек Кунгужинов «Викинг» фильміндегі жұмысы үшін Taurus World Stunt Awards сыйлығына ұсынылды.'],
 ['2018','Жайдарбек Кунгужинов пен Nomad Stunts тобы Қазақстанның КИНЭС рекордтар кітабына енді. Жайдарбек Кунгужинов сондай-ақ Қазақстан Республикасының мемлекеттік наградасы — «Құрмет» орденінің иегері атанды.'],
 ['2020','Жайдарбек Кунгужинов Әлемдік каскадерлер академиясының (АҚШ) «Таурус» сыйлығын алған алғашқы қазақстандық болып, тарих жасады. Бұл беделді марапатты каскадерлердің «Оскары» деп атайды. Сыйлық «Томирис» фильміндегі трюктерді қойғаны үшін берілді.'],
 ['Бүгін','Топ жұмыс істеген жылдар ішінде Nomad Stunts каскадерлері Киану Ривз, Колин Фаррелл, Хоакин Феникс, Мишель Йео, Сирша Ронан, Сэмюэл Л. Джексон, Натали Портман, Лю Ифэй, Шон Бин, Эд Харрис, Брюс Уиллис, Джейсон Стейтем, Таданобу Асано, Арнольд Шварценеггер, Джет Ли, Донни Йен, Том Ву, Бенедикт Вонг сияқты әлем киносының жұлдыздарымен жұмыс істеді. Бүгінде Nomad Stunts — Қазақстан, Қырғызстан, Австралия, Канада, Түркия, Әзірбайжан, Ресей және Украинадан келген 86 каскадер, 40 үйретілген трюк жылқысы және Алматы облысы, Есік қаласының маңындағы арнайы жаттығу базасы.']]
};

const STARS = {
 en:['Keanu Reeves','Colin Farrell','Joaquin Phoenix','Michelle Yeoh','Saoirse Ronan','Samuel L. Jackson','Natalie Portman','Liu Yifei','Sean Bean','Ed Harris','Bruce Willis','Jason Statham','Tadanobu Asano','Arnold Schwarzenegger','Jet Li','Donnie Yen','Tom Wu','Benedict Wong'],
 ru:['Киану Ривз','Колин Фаррелл','Хоакин Феникс','Мишель Йео','Сирша Ронан','Сэмюэл Л. Джексон','Натали Портман','Лю Ифэй','Шон Бин','Эд Харрис','Брюс Уиллис','Джейсон Стейтем','Таданобу Асано','Арнольд Шварценеггер','Джет Ли','Донни Йен','Том Ву','Бенедикт Вонг']
};
STARS.kk = STARS.ru;

const FEATURED = ['napoleon','rrr','mulan','marco-polo','viking','tomiris','neuderzhimye-2','47-roninov','conan-the-barbarian','the-letter-for-the-king','mcmafia','thugs-of-hindostan','baahubali-rozhdenie-legendy','awakening-the-great-seljuk','the-way-back','your-highness','fatih'];

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
const NAV = [['about','about.html'],['team','team.html'],['projects','projects.html'],['gallery','gallery.html'],['contacts','contacts.html']];
const soc = () => SOCIAL.map(([k,n,u])=>`<a href="${u}" target="_blank" rel="noopener" aria-label="${n}"><svg viewBox="0 0 24 24">${ICONS[k]}</svg></a>`).join('');

function chrome(){
  const hdr = document.createElement('header');
  hdr.innerHTML = `<a href="index.html" class="logo"><img src="assets/logo.svg" alt="Nomad Stunts" width="167" height="76"></a>
  <nav id="nav">${NAV.map(([k,u])=>`<a href="${u}" class="${k===PAGE?'on':''}">${t('nav_'+k)}</a>`).join('')}</nav>
  <div class="hdr-r"><div class="soc">${soc()}</div>
  <div class="lang" role="group" aria-label="Language">${LANGS.map(([c,n])=>`<button data-l="${c}" class="${c===lang?'on':''}" lang="${c}">${n}</button>`).join('')}</div></div>
  <button class="burger" aria-label="Menu"><i></i></button>`;
  const old = $('header'); old ? old.replaceWith(hdr) : document.body.prepend(hdr);

  const ft = document.createElement('footer'); ft.id='footer';
  ft.innerHTML = `<div class="wrap ft">
   <div><a href="index.html" class="ft-logo"><img src="assets/logo.svg" alt="Nomad Stunts" width="167" height="76"></a><h4>${t('ft_menu')}</h4><ul><li><a href="index.html">${t('home')}</a></li>${NAV.map(([k,u])=>`<li><a href="${u}">${t('nav_'+k)}</a></li>`).join('')}</ul></div>
   <div><h4>${t('ft_contact')}</h4><p>${t('ft_p')}</p><p class="agent">${t('agent')}</p>
    <div class="ph"><a href="tel:+77779900068">+7 777 99 000 68</a><a href="tel:+77770345034">+7 777 03 450 34</a></div>
    <p style="margin-top:8px">Email: <a href="mailto:nomadstunts@gmail.com">nomadstunts@gmail.com</a></p></div>
   <div class="ft-r"><div class="soc">${soc()}</div><a class="btn" href="mailto:nomadstunts@gmail.com">${t('send')}</a></div>
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
function showreel(){
  const src = SHOWREEL_MP4
    ? `<video src="${SHOWREEL_MP4}" controls autoplay playsinline style="width:100%;height:100%"></video>`
    : `<iframe src="https://www.youtube-nocookie.com/embed/${SHOWREEL}?autoplay=1&rel=0" allow="autoplay; encrypted-media; fullscreen" allowfullscreen title="Nomad Stunts Showreel"></iframe>`;
  openModal('vid', `<div class="vwrap">${src}</div>`);
}

// ---------- people & films ----------
const fullName = p => L(p.n).join(' ');
const roles = p => p.roles.map(r=>t('r_'+r)).join(', ');
function personModal(p, isFounder){
  const sk = p.sk.map(([k,v])=>`<div class="bar"><span>${esc(L(D.skills[k])||k)}<b>${v}%</b></span><i style="--w:${v}%"></i></div>`).join('');
  const hon = L(p.h);
  const role = isFounder ? L(p.role) : roles(p);
  openModal('person', `<div class="mimg"><img src="${img(p.photo||p.img)}" alt="${esc(fullName(p))}"></div>
   <div class="mtxt"><div class="eyebrow">${isFounder ? t('cat_founder') : p.cats.map(c=>t('cat_'+c)).join(' · ')}</div>
   <h2 class="h">${esc(L(p.n)[0])}<br>${esc(L(p.n)[1])}</h2><p class="role">${esc(role)}</p>
   ${sk ? `<h4>${t('skills')}</h4><div class="bars">${sk}</div>` : ''}
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
function orderedFilms(){ const fe = FEATURED.map(id=>D.films.find(f=>f.id===id)).filter(Boolean); return [...fe, ...D.films.filter(f=>!fe.includes(f))]; }

function reel(el){
  el.innerHTML = `<button class="arr p" aria-label="‹">&#8592;</button><div class="track">${orderedFilms().map(posterCard).join('')}</div><button class="arr n" aria-label="›">&#8594;</button>`;
  const tr = $('.track',el);
  $('.arr.n',el).onclick=()=>tr.scrollBy({left:tr.clientWidth*.8,behavior:'smooth'});
  $('.arr.p',el).onclick=()=>tr.scrollBy({left:-tr.clientWidth*.8,behavior:'smooth'});
  bindPosters(tr);
}
function marquee(el){ const s = STARS[lang].map(x=>`${esc(x)}<i>&#9670;</i>`).join(''); el.innerHTML = `<div>${s+s}</div>`; el.setAttribute('aria-label',t('worked')); }
function serviceCards(el){ el.innerHTML = SERVICES.map(s=>`<article class="card rv"><h3>${s[lang][0]}</h3><div class="ph"><img src="${img(s.img)}" alt="" loading="lazy"></div><p>${s[lang][1]}</p></article>`).join(''); }
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
    const HERO = ['big/44798_710250_00.jpg','21109_783683_14.jpg','big/97241_93672_00.jpg'];
    HERO.forEach((h,i)=>{ const d=document.createElement('div'); d.className='sl'+(i?'':' on'); d.style.backgroundImage=`url(${img(h)})`; box.prepend(d);
      const b=document.createElement('button'); b.setAttribute('aria-label','Slide '+(i+1)); if(!i) b.className='on'; b.onclick=()=>go(i); $('#dots').append(b); });
    const slides = $$('.sl',box).reverse(); let cur=0;
    const go = i => { slides[cur].classList.remove('on'); $('#dots').children[cur].classList.remove('on'); cur=i; slides[cur].classList.add('on'); $('#dots').children[cur].classList.add('on'); clearInterval(heroTimer); heroTimer=setInterval(()=>go((cur+1)%slides.length),6000); };
    go(0);
    // muted background showreel over the slideshow
    if (!matchMedia('(prefers-reduced-motion: reduce)').matches && (SHOWREEL_MP4 || location.protocol!=='file:')){
      const v = document.createElement('div'); v.className='vid';
      v.innerHTML = SHOWREEL_MP4
        ? `<video src="${SHOWREEL_MP4}" autoplay muted loop playsinline style="width:100%;height:100%;object-fit:cover"></video>`
        : `<iframe src="https://www.youtube-nocookie.com/embed/${SHOWREEL}?autoplay=1&mute=1&loop=1&playlist=${SHOWREEL}&controls=0&rel=0&playsinline=1&iv_load_policy=3&disablekb=1" allow="autoplay; encrypted-media" tabindex="-1" aria-hidden="true" title="Showreel"></iframe>`;
      $('.dots',box).before(v);
      const show = () => setTimeout(()=>{ v.classList.add('on'); clearInterval(heroTimer); $('#dots').style.display='none'; }, SHOWREEL_MP4?0:1800);
      (v.firstElementChild).addEventListener(SHOWREEL_MP4?'playing':'load', show, {once:true});
    }
    $('#play').onclick = showreel;
  }
  stats($('#stats')); serviceCards($('#cards')); founderBlock($('#founder'));
  const tabs = ['coord','horse','rig'];
  const cur = $('#tabs').dataset.t || 'coord';
  $('#tabs').innerHTML = tabs.map(k=>`<button data-t="${k}" class="${k===cur?'on':''}">${t('cat_'+k)}</button>`).join('');
  $('#tabs').onclick = e => { const b=e.target.closest('button'); if(b){ $('#tabs').dataset.t=b.dataset.t; PAGES.home(); } };
  const list = D.team.filter(p=>p.cats.includes(cur));
  $('#grid').innerHTML = list.map((p)=>memCard(p,D.team.indexOf(p))).join('');
  $('#grid').onclick = e => { const b=e.target.closest('[data-p]'); if(b) personModal(D.team[b.dataset.p]); };
  reel($('#reel')); marquee($('#stars'));
 },
 about(){
  timeline($('#tl'));
  $('#dirs').innerHTML = DIRECTIONS.map(d=>`<div class="rv"><img src="${img(d[0])}" alt="" loading="lazy"><span>${esc({en:d[1],ru:d[2],kk:d[3]}[lang])}</span></div>`).join('');
  founderBlock($('#founder')); stats($('#stats')); marquee($('#stars'));
  const v=$('#video'); v.style.backgroundImage='url(assets/video/showreel.jpg)'; v.onclick=showreel;
  $('#video .lbl').textContent = t('watch');
 },
 team(){
  const st = $('#tabs').dataset; const cur = st.t || 'all'; const q = ($('#q').value||'').trim().toLowerCase();
  const cats = ['coord','horse','rig','perf'];
  $('#tabs').innerHTML = `<button data-t="all" class="${cur==='all'?'on':''}">${t('all')}<sup>${D.team.length+1}</sup></button><button data-t="founder">${t('cat_founder')}</button>` +
    cats.map(k=>`<button data-t="${k}" class="${k===cur?'on':''}">${t('cat_'+k)}<sup>${D.team.filter(p=>p.cats.includes(k)).length}</sup></button>`).join('');
  $('#tabs').onclick = e => { const b=e.target.closest('button'); if(!b) return; if (b.dataset.t==='founder') return personModal(D.founder,true); st.t=b.dataset.t; PAGES.team(); };
  $('#q').placeholder = t('search'); $('#q').oninput = PAGES.team;
  const match = p => !q || ['en','ru'].some(l=>p.n[l].join(' ').toLowerCase().includes(q));
  let list = D.team.filter(p=>(cur==='all'||p.cats.includes(cur)) && match(p));
  const showF = (cur==='all') && match(D.founder);
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
  if (!g.children.length) g.innerHTML = D.gallery.map((x,i)=>`<button data-g="${i}" class="rv"><img src="${UP}small/${x}" data-full="${UP+x}" alt="Nomad Stunts ${i+1}" loading="lazy" onerror="this.onerror=null;this.src=this.dataset.full"></button>`).join('');
  g.onclick = e => { const b=e.target.closest('[data-g]'); if (b) lightbox(+b.dataset.g); };
 },
 contacts(){}
};
function lightbox(i){
  const n = D.gallery.length;
  const show = j => { i=(j+n)%n; $('img',modal).src = UP + D.gallery[i]; };
  openModal('lb', `<img alt=""><button class="arr p" aria-label="‹">&#8592;</button><button class="arr n" aria-label="›">&#8594;</button>`, {ArrowLeft:()=>show(i-1), ArrowRight:()=>show(i+1)});
  $('.arr.p',modal).onclick = e => { e.stopPropagation(); show(i-1); };
  $('.arr.n',modal).onclick = e => { e.stopPropagation(); show(i+1); };
  show(i);
}

function render(){
  chrome(); i18nStatic();
  document.title = (PAGE==='home' ? 'Nomad Stunts | ' + t('tagline') : t('nav_'+PAGE) + ' | Nomad Stunts');
  (PAGES[PAGE]||(()=>{}))();
  watch();
}
document.addEventListener('click', e => { if (e.target.closest('nav a')) document.body.classList.remove('menu'); });
render();
})();
