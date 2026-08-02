// Copyright 2026 Rangertaha. All Rights Reserved.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/rangertaha/urlinsane/internal/plugins/languages"
	_ "github.com/rangertaha/urlinsane/internal/plugins/languages/all"
	"github.com/urfave/cli/v2"
)

// defaultSynonyms returns brand-adjacent word groups per language -- the nouns and
// verbs an attacker most often bolts onto a target name (shop, login, secure, pay...).
func defaultSynonyms(langID string) [][]string {
	switch langID {
	case "en":
		return [][]string{
			{"begin", "start", "commence", "initiate"},
			{"shop", "store", "market", "buy"},
			{"login", "signin", "logon", "account"},
			{"secure", "safe", "protected", "trusted"},
			{"support", "help", "service", "assist"},
			{"pay", "payment", "checkout", "billing"},
			{"bank", "finance", "money", "wallet"},
			{"mail", "email", "message", "inbox"},
			{"cloud", "server", "host", "hosting"},
			{"web", "site", "online", "net", "portal"},
			{"news", "blog", "press", "media"},
			{"free", "gift", "bonus", "promo"},
			{"best", "top", "premium", "pro"},
			{"new", "latest", "fresh", "modern"},
			{"fast", "quick", "rapid", "speed"},
			{"app", "application", "software", "tool"},
			{"download", "get", "install"},
			{"contact", "reach", "connect"},
			{"search", "find", "lookup"},
			{"update", "upgrade", "renew"},
		}
	case "fr":
		return [][]string{
			{"debut", "commencement", "depart"},
			{"boutique", "magasin", "marche", "achat"},
			{"connexion", "identifiant", "compte"},
			{"securise", "sur", "protege", "fiable"},
			{"support", "aide", "service", "assistance"},
			{"paiement", "reglement", "facturation"},
			{"banque", "finance", "argent", "portefeuille"},
			{"courrier", "email", "message", "boite"},
			{"nuage", "serveur", "hebergement"},
			{"web", "site", "enligne", "reseau", "portail"},
			{"actualites", "blog", "presse", "medias"},
			{"gratuit", "cadeau", "bonus", "promo"},
			{"meilleur", "top", "premium", "pro"},
			{"nouveau", "dernier", "recent", "moderne"},
			{"rapide", "vite", "express"},
			{"appli", "application", "logiciel", "outil"},
			{"telecharger", "obtenir", "installer"},
			{"contact", "joindre", "connecter"},
			{"recherche", "trouver"},
			{"mise", "a", "jour", "renouveler"},
		}
	case "es":
		return [][]string{
			{"inicio", "comienzo", "arranque"},
			{"tienda", "comercio", "mercado", "compra"},
			{"acceso", "ingreso", "cuenta"},
			{"seguro", "protegido", "confiable"},
			{"soporte", "ayuda", "servicio", "asistencia"},
			{"pago", "cobro", "facturacion"},
			{"banco", "finanzas", "dinero", "cartera"},
			{"correo", "email", "mensaje", "bandeja"},
			{"nube", "servidor", "alojamiento"},
			{"web", "sitio", "enlinea", "red", "portal"},
			{"noticias", "blog", "prensa", "medios"},
			{"gratis", "regalo", "bono", "promo"},
			{"mejor", "top", "premium", "pro"},
			{"nuevo", "ultimo", "reciente", "moderno"},
			{"rapido", "veloz", "express"},
			{"app", "aplicacion", "software", "herramienta"},
			{"descargar", "obtener", "instalar"},
			{"contacto", "conectar"},
			{"buscar", "encontrar"},
			{"actualizar", "renovar"},
		}
	case "pt":
		return [][]string{
			{"inicio", "comeco", "arranque"},
			{"loja", "comercio", "mercado", "compra"},
			{"acesso", "entrar", "conta"},
			{"seguro", "protegido", "confiavel"},
			{"suporte", "ajuda", "servico", "assistencia"},
			{"pagamento", "cobranca", "faturamento"},
			{"banco", "financas", "dinheiro", "carteira"},
			{"correio", "email", "mensagem", "caixa"},
			{"nuvem", "servidor", "hospedagem"},
			{"web", "site", "online", "rede", "portal"},
			{"noticias", "blog", "imprensa", "midia"},
			{"gratis", "presente", "bonus", "promo"},
			{"melhor", "top", "premium", "pro"},
			{"novo", "ultimo", "recente", "moderno"},
			{"rapido", "veloz", "express"},
			{"app", "aplicativo", "software", "ferramenta"},
			{"baixar", "obter", "instalar"},
			{"contato", "conectar"},
			{"buscar", "encontrar"},
			{"atualizar", "renovar"},
		}
	case "it":
		return [][]string{
			{"inizio", "partenza", "avvio"},
			{"negozio", "bottega", "mercato", "acquisto"},
			{"accesso", "login", "conto"},
			{"sicuro", "protetto", "affidabile"},
			{"supporto", "aiuto", "servizio", "assistenza"},
			{"pagamento", "saldo", "fatturazione"},
			{"banca", "finanza", "denaro", "portafoglio"},
			{"posta", "email", "messaggio", "casella"},
			{"nuvola", "server", "hosting"},
			{"web", "sito", "online", "rete", "portale"},
			{"notizie", "blog", "stampa", "media"},
			{"gratis", "regalo", "bonus", "promo"},
			{"migliore", "top", "premium", "pro"},
			{"nuovo", "ultimo", "recente", "moderno"},
			{"rapido", "veloce", "espresso"},
			{"app", "applicazione", "software", "strumento"},
			{"scaricare", "ottenere", "installare"},
			{"contatto", "connettere"},
			{"cerca", "trova"},
			{"aggiorna", "rinnova"},
		}
	case "de":
		return [][]string{
			{"anfang", "beginn", "start"},
			{"shop", "laden", "markt", "kauf"},
			{"anmeldung", "login", "konto"},
			{"sicher", "geschuetzt", "vertrauenswuerdig"},
			{"support", "hilfe", "service", "beistand"},
			{"zahlung", "bezahlung", "abrechnung"},
			{"bank", "finanzen", "geld", "brieftasche"},
			{"post", "email", "nachricht", "postfach"},
			{"cloud", "server", "hosting"},
			{"web", "seite", "online", "netz", "portal"},
			{"nachrichten", "blog", "presse", "medien"},
			{"gratis", "geschenk", "bonus", "aktion"},
			{"beste", "top", "premium", "pro"},
			{"neu", "neueste", "aktuell", "modern"},
			{"schnell", "rasch", "express"},
			{"app", "anwendung", "software", "werkzeug"},
			{"herunterladen", "holen", "installieren"},
			{"kontakt", "verbinden"},
			{"suche", "finden"},
			{"aktualisieren", "erneuern"},
		}
	case "nl":
		return [][]string{
			{"begin", "start", "aanvang"},
			{"winkel", "shop", "markt", "koop"},
			{"inloggen", "aanmelden", "account"},
			{"veilig", "beschermd", "betrouwbaar"},
			{"support", "hulp", "service", "bijstand"},
			{"betaling", "afrekenen", "facturering"},
			{"bank", "financien", "geld", "portemonnee"},
			{"post", "email", "bericht", "inbox"},
			{"cloud", "server", "hosting"},
			{"web", "site", "online", "netwerk", "portaal"},
			{"nieuws", "blog", "pers", "media"},
			{"gratis", "cadeau", "bonus", "actie"},
			{"beste", "top", "premium", "pro"},
			{"nieuw", "laatste", "recent", "modern"},
			{"snel", "vlug", "express"},
			{"app", "applicatie", "software", "tool"},
			{"downloaden", "ophalen", "installeren"},
			{"contact", "verbinden"},
			{"zoeken", "vinden"},
			{"bijwerken", "vernieuwen"},
		}
	case "sv":
		return [][]string{
			{"start", "borja", "begynnelse"},
			{"butik", "affar", "marknad", "kop"},
			{"inloggning", "logga", "konto"},
			{"saker", "skyddad", "palitlig"},
			{"support", "hjalp", "service", "bistand"},
			{"betalning", "kassa", "fakturering"},
			{"bank", "finans", "pengar", "planbok"},
			{"post", "email", "meddelande", "inkorg"},
			{"moln", "server", "webbhotell"},
			{"webb", "sida", "online", "natverk", "portal"},
			{"nyheter", "blogg", "press", "media"},
			{"gratis", "present", "bonus", "kampanj"},
			{"basta", "topp", "premium", "pro"},
			{"ny", "senaste", "aktuell", "modern"},
			{"snabb", "kvick", "express"},
			{"app", "applikation", "programvara", "verktyg"},
			{"ladda", "hamta", "installera"},
			{"kontakt", "ansluta"},
			{"sok", "hitta"},
			{"uppdatera", "fornya"},
		}
	case "no":
		return [][]string{
			{"start", "begynnelse", "begynne"},
			{"butikk", "handel", "marked", "kjop"},
			{"innlogging", "logg", "konto"},
			{"sikker", "beskyttet", "palitelig"},
			{"support", "hjelp", "service", "bistand"},
			{"betaling", "kasse", "fakturering"},
			{"bank", "finans", "penger", "lommebok"},
			{"post", "email", "melding", "innboks"},
			{"sky", "server", "webhotell"},
			{"web", "side", "online", "nettverk", "portal"},
			{"nyheter", "blogg", "presse", "medier"},
			{"gratis", "gave", "bonus", "kampanje"},
			{"beste", "topp", "premium", "pro"},
			{"ny", "nyeste", "aktuell", "moderne"},
			{"rask", "kjapp", "ekspress"},
			{"app", "applikasjon", "programvare", "verktoy"},
			{"laste", "hente", "installere"},
			{"kontakt", "koble"},
			{"sok", "finn"},
			{"oppdater", "forny"},
		}
	case "da":
		return [][]string{
			{"start", "begyndelse", "begynde"},
			{"butik", "handel", "marked", "kob"},
			{"login", "logon", "konto"},
			{"sikker", "beskyttet", "palidelig"},
			{"support", "hjaelp", "service", "bistand"},
			{"betaling", "kasse", "fakturering"},
			{"bank", "finans", "penge", "pung"},
			{"post", "email", "besked", "indbakke"},
			{"sky", "server", "webhotel"},
			{"web", "side", "online", "netvaerk", "portal"},
			{"nyheder", "blog", "presse", "medier"},
			{"gratis", "gave", "bonus", "kampagne"},
			{"bedste", "top", "premium", "pro"},
			{"ny", "nyeste", "aktuel", "moderne"},
			{"hurtig", "kvik", "ekspres"},
			{"app", "applikation", "software", "vaerktoj"},
			{"hente", "downloade", "installere"},
			{"kontakt", "forbinde"},
			{"sog", "find"},
			{"opdater", "forny"},
		}
	case "fi":
		return [][]string{
			{"alku", "aloitus"},
			{"kauppa", "myymala", "markkina", "osto"},
			{"kirjautuminen", "tunnus", "tili"},
			{"turvallinen", "suojattu", "luotettava"},
			{"tuki", "apu", "palvelu", "avustus"},
			{"maksu", "maksaminen", "laskutus"},
			{"pankki", "rahoitus", "raha", "lompakko"},
			{"posti", "sahkoposti", "viesti"},
			{"pilvi", "palvelin", "webhotelli"},
			{"web", "sivusto", "verkossa", "verkko", "portaali"},
			{"uutiset", "blogi", "lehdisto", "media"},
			{"ilmainen", "lahja", "bonus", "tarjous"},
			{"paras", "huippu", "premium", "pro"},
			{"uusi", "uusin", "ajankohtainen", "moderni"},
			{"nopea", "vikkela", "pika"},
			{"sovellus", "ohjelma", "tyokalu"},
			{"lataa", "hanki", "asenna"},
			{"yhteys", "yhdista"},
			{"haku", "etsi"},
			{"paivita", "uusi"},
		}
	case "ru":
		return [][]string{
			{"начало", "старт"},
			{"магазин", "лавка", "рынок", "покупка"},
			{"вход", "логин", "аккаунт"},
			{"безопасный", "защищенный", "надежный"},
			{"поддержка", "помощь", "сервис"},
			{"оплата", "платеж", "расчет"},
			{"банк", "финансы", "деньги", "кошелек"},
			{"почта", "email", "сообщение"},
			{"облако", "сервер", "хостинг"},
			{"веб", "сайт", "онлайн", "сеть", "портал"},
			{"новости", "блог", "пресса", "медиа"},
			{"бесплатно", "подарок", "бонус", "акция"},
			{"лучший", "топ", "премиум", "про"},
			{"новый", "последний", "свежий", "современный"},
			{"быстрый", "скорый", "экспресс"},
			{"приложение", "программа", "инструмент"},
			{"скачать", "получить", "установить"},
			{"контакт", "связаться"},
			{"поиск", "найти"},
			{"обновить", "продлить"},
		}
	case "uk":
		return [][]string{
			{"початок", "старт"},
			{"магазин", "крамниця", "ринок", "покупка"},
			{"вхід", "логін", "акаунт"},
			{"безпечний", "захищений", "надійний"},
			{"підтримка", "допомога", "сервіс"},
			{"оплата", "платіж", "розрахунок"},
			{"банк", "фінанси", "гроші", "гаманець"},
			{"пошта", "email", "повідомлення"},
			{"хмара", "сервер", "хостинг"},
			{"веб", "сайт", "онлайн", "мережа", "портал"},
			{"новини", "блог", "преса", "медіа"},
			{"безкоштовно", "подарунок", "бонус", "акція"},
			{"кращий", "топ", "преміум", "про"},
			{"новий", "останній", "свіжий", "сучасний"},
			{"швидкий", "хуткий", "експрес"},
			{"додаток", "програма", "інструмент"},
			{"завантажити", "отримати", "встановити"},
			{"контакт", "звязатися"},
			{"пошук", "знайти"},
			{"оновити", "продовжити"},
		}
	case "pl":
		return [][]string{
			{"poczatek", "start"},
			{"sklep", "handel", "rynek", "zakup"},
			{"logowanie", "konto"},
			{"bezpieczny", "chroniony", "zaufany"},
			{"wsparcie", "pomoc", "serwis"},
			{"platnosc", "zaplata", "rozliczenie"},
			{"bank", "finanse", "pieniadze", "portfel"},
			{"poczta", "email", "wiadomosc"},
			{"chmura", "serwer", "hosting"},
			{"web", "strona", "online", "siec", "portal"},
			{"wiadomosci", "blog", "prasa", "media"},
			{"darmowy", "prezent", "bonus", "promocja"},
			{"najlepszy", "top", "premium", "pro"},
			{"nowy", "najnowszy", "swiezy", "nowoczesny"},
			{"szybki", "predki", "ekspres"},
			{"aplikacja", "program", "narzedzie"},
			{"pobierz", "uzyskaj", "zainstaluj"},
			{"kontakt", "polacz"},
			{"szukaj", "znajdz"},
			{"aktualizuj", "odnow"},
		}
	case "cs":
		return [][]string{
			{"zacatek", "start"},
			{"obchod", "trh", "nakup"},
			{"prihlaseni", "ucet"},
			{"bezpecny", "chraneny", "duveryhodny"},
			{"podpora", "pomoc", "servis"},
			{"platba", "uhrada", "vyuctovani"},
			{"banka", "finance", "penize", "penezenka"},
			{"posta", "email", "zprava"},
			{"cloud", "server", "hosting"},
			{"web", "stranka", "online", "sit", "portal"},
			{"zpravy", "blog", "tisk", "media"},
			{"zdarma", "darek", "bonus", "akce"},
			{"nejlepsi", "top", "premium", "pro"},
			{"novy", "nejnovejsi", "cerstvy", "moderni"},
			{"rychly", "svizny", "expres"},
			{"aplikace", "program", "nastroj"},
			{"stahnout", "ziskat", "nainstalovat"},
			{"kontakt", "spojit"},
			{"hledat", "najit"},
			{"aktualizovat", "obnovit"},
		}
	case "tr":
		return [][]string{
			{"baslangic", "start"},
			{"magaza", "dukkan", "pazar", "alisveris"},
			{"giris", "oturum", "hesap"},
			{"guvenli", "korumali", "guvenilir"},
			{"destek", "yardim", "servis"},
			{"odeme", "tahsilat", "faturalama"},
			{"banka", "finans", "para", "cuzdan"},
			{"posta", "email", "mesaj"},
			{"bulut", "sunucu", "barindirma"},
			{"web", "site", "online", "ag", "portal"},
			{"haber", "blog", "basin", "medya"},
			{"ucretsiz", "hediye", "bonus", "kampanya"},
			{"en", "iyi", "top", "premium", "pro"},
			{"yeni", "son", "guncel", "modern"},
			{"hizli", "cabuk", "ekspres"},
			{"uygulama", "program", "arac"},
			{"indir", "edin", "kur"},
			{"iletisim", "baglan"},
			{"ara", "bul"},
			{"guncelle", "yenile"},
		}
	case "el":
		return [][]string{
			{"αρχη", "εναρξη"},
			{"καταστημα", "μαγαζι", "αγορα"},
			{"συνδεση", "λογαριασμος"},
			{"ασφαλης", "προστατευμενος", "αξιοπιστος"},
			{"υποστηριξη", "βοηθεια", "υπηρεσια"},
			{"πληρωμη", "εξοφληση", "χρεωση"},
			{"τραπεζα", "οικονομικα", "χρηματα", "πορτοφολι"},
			{"ταχυδρομειο", "email", "μηνυμα"},
			{"συννεφο", "διακομιστης", "φιλοξενια"},
			{"web", "ιστοτοπος", "online", "δικτυο", "πυλη"},
			{"νεα", "ιστολογιο", "τυπος", "μεσα"},
			{"δωρεαν", "δωρο", "μπονους", "προσφορα"},
			{"καλυτερο", "κορυφαιο", "premium", "pro"},
			{"νεο", "τελευταιο", "προσφατο", "μοντερνο"},
			{"γρηγορο", "ταχυ", "εξπρες"},
			{"εφαρμογη", "προγραμμα", "εργαλειο"},
			{"καταβασε", "παρε", "εγκατεστησε"},
			{"επαφη", "συνδεση"},
			{"αναζητηση", "βρες"},
			{"ενημερωση", "ανανεωση"},
		}
	case "ar":
		return [][]string{
			{"بداية", "بدء", "انطلاق"},
			{"متجر", "محل", "سوق", "شراء"},
			{"دخول", "تسجيل", "حساب"},
			{"امن", "محمي", "موثوق"},
			{"دعم", "مساعدة", "خدمة"},
			{"دفع", "سداد", "فوترة"},
			{"بنك", "مالية", "نقود", "محفظة"},
			{"بريد", "ايميل", "رسالة"},
			{"سحابة", "خادم", "استضافة"},
			{"ويب", "موقع", "اونلاين", "شبكة", "بوابة"},
			{"اخبار", "مدونة", "صحافة", "اعلام"},
			{"مجاني", "هدية", "مكافأة", "عرض"},
			{"افضل", "قمة", "مميز"},
			{"جديد", "احدث", "حديث", "عصري"},
			{"سريع", "عاجل", "اكسبرس"},
			{"تطبيق", "برنامج", "اداة"},
			{"تحميل", "تنزيل", "تثبيت"},
			{"اتصال", "تواصل"},
			{"بحث", "ايجاد"},
			{"تحديث", "تجديد"},
		}
	case "fa":
		return [][]string{
			{"شروع", "آغاز"},
			{"فروشگاه", "مغازه", "بازار", "خرید"},
			{"ورود", "لاگین", "حساب"},
			{"امن", "محافظت", "مطمئن"},
			{"پشتیبانی", "کمک", "خدمات"},
			{"پرداخت", "تسویه", "صورتحساب"},
			{"بانک", "مالی", "پول", "کیف"},
			{"پست", "ایمیل", "پیام"},
			{"ابر", "سرور", "میزبانی"},
			{"وب", "سایت", "آنلاین", "شبکه", "درگاه"},
			{"اخبار", "بلاگ", "مطبوعات", "رسانه"},
			{"رایگان", "هدیه", "پاداش", "تخفیف"},
			{"بهترین", "برتر", "ویژه"},
			{"جدید", "تازه", "نو", "مدرن"},
			{"سریع", "تند", "اکسپرس"},
			{"اپلیکیشن", "برنامه", "ابزار"},
			{"دانلود", "دریافت", "نصب"},
			{"تماس", "ارتباط"},
			{"جستجو", "یافتن"},
			{"بروزرسانی", "تمدید"},
		}
	case "iw":
		return [][]string{
			{"התחלה", "תחילה"},
			{"חנות", "מסחר", "שוק", "קניה"},
			{"כניסה", "התחברות", "חשבון"},
			{"בטוח", "מוגן", "אמין"},
			{"תמיכה", "עזרה", "שירות"},
			{"תשלום", "סליקה", "חיוב"},
			{"בנק", "פיננסים", "כסף", "ארנק"},
			{"דואר", "אימייל", "הודעה"},
			{"ענן", "שרת", "אחסון"},
			{"אתר", "אינטרנט", "מקוון", "רשת", "פורטל"},
			{"חדשות", "בלוג", "עיתונות", "מדיה"},
			{"חינם", "מתנה", "בונוס", "מבצע"},
			{"הטוב", "מוביל", "פרימיום"},
			{"חדש", "אחרון", "עדכני", "מודרני"},
			{"מהיר", "זריז", "אקספרס"},
			{"אפליקציה", "תוכנה", "כלי"},
			{"הורדה", "קבלה", "התקנה"},
			{"קשר", "יצירת"},
			{"חיפוש", "מציאה"},
			{"עדכון", "חידוש"},
		}
	case "hi":
		return [][]string{
			{"शुरुआत", "आरंभ", "प्रारंभ"},
			{"दुकान", "बाजार", "खरीद"},
			{"लॉगिन", "प्रवेश", "खाता"},
			{"सुरक्षित", "संरक्षित", "विश्वसनीय"},
			{"सहायता", "मदद", "सेवा"},
			{"भुगतान", "अदायगी", "बिलिंग"},
			{"बैंक", "वित्त", "पैसा", "बटुआ"},
			{"डाक", "ईमेल", "संदेश"},
			{"क्लाउड", "सर्वर", "होस्टिंग"},
			{"वेब", "साइट", "ऑनलाइन", "नेटवर्क", "पोर्टल"},
			{"समाचार", "ब्लॉग", "प्रेस", "मीडिया"},
			{"मुफ्त", "उपहार", "बोनस", "ऑफर"},
			{"सर्वोत्तम", "शीर्ष", "प्रीमियम"},
			{"नया", "नवीनतम", "ताजा", "आधुनिक"},
			{"तेज", "शीघ्र", "एक्सप्रेस"},
			{"ऐप", "एप्लिकेशन", "सॉफ्टवेयर", "उपकरण"},
			{"डाउनलोड", "प्राप्त", "इंस्टॉल"},
			{"संपर्क", "जुड़ना"},
			{"खोज", "ढूंढ"},
			{"अपडेट", "नवीनीकरण"},
		}
	case "zh":
		return [][]string{
			{"开始", "开端", "启动"},
			{"商店", "商城", "市场", "购买"},
			{"登录", "登入", "账户"},
			{"安全", "保护", "可信"},
			{"支持", "帮助", "服务"},
			{"支付", "付款", "结算"},
			{"银行", "金融", "钱", "钱包"},
			{"邮件", "邮箱", "消息"},
			{"云", "服务器", "主机"},
			{"网", "网站", "在线", "网络", "门户"},
			{"新闻", "博客", "媒体"},
			{"免费", "礼物", "奖金", "优惠"},
			{"最佳", "顶级", "高级"},
			{"新", "最新", "现代"},
			{"快速", "迅速", "快捷"},
			{"应用", "程序", "软件", "工具"},
			{"下载", "获取", "安装"},
			{"联系", "连接"},
			{"搜索", "查找"},
			{"更新", "续期"},
		}
	case "ja":
		return [][]string{
			{"開始", "始まり", "起動"},
			{"店舗", "ショップ", "市場", "購入"},
			{"ログイン", "サインイン", "アカウント"},
			{"安全", "保護", "信頼"},
			{"サポート", "ヘルプ", "サービス"},
			{"支払", "決済", "請求"},
			{"銀行", "金融", "お金", "財布"},
			{"メール", "郵便", "メッセージ"},
			{"クラウド", "サーバー", "ホスティング"},
			{"ウェブ", "サイト", "オンライン", "ネット"},
			{"ニュース", "ブログ", "報道", "メディア"},
			{"無料", "ギフト", "ボーナス"},
			{"最高", "トップ", "プレミアム"},
			{"新しい", "最新", "現代"},
			{"高速", "迅速", "エクスプレス"},
			{"アプリ", "アプリケーション", "ソフト"},
			{"ダウンロード", "取得", "インストール"},
			{"連絡", "接続"},
			{"検索", "探す"},
			{"更新", "更改"},
		}
	case "ko":
		return [][]string{
			{"시작", "개시", "출발"},
			{"상점", "쇼핑", "시장", "구매"},
			{"로그인", "접속", "계정"},
			{"안전", "보호", "신뢰"},
			{"지원", "도움", "서비스"},
			{"결제", "지불", "청구"},
			{"은행", "금융", "돈", "지갑"},
			{"메일", "우편", "메시지"},
			{"클라우드", "서버", "호스팅"},
			{"웹", "사이트", "온라인", "네트워크"},
			{"뉴스", "블로그", "언론", "미디어"},
			{"무료", "선물", "보너스", "혜택"},
			{"최고", "최상", "프리미엄"},
			{"새로운", "최신", "현대"},
			{"빠른", "신속", "익스프레스"},
			{"앱", "애플리케이션", "소프트웨어"},
			{"다운로드", "받기", "설치"},
			{"연락", "연결"},
			{"검색", "찾기"},
			{"업데이트", "갱신"},
		}
	case "th":
		return [][]string{
			{"เริ่ม", "เริ่มต้น"},
			{"ร้าน", "ร้านค้า", "ตลาด", "ซื้อ"},
			{"เข้าสู่ระบบ", "ล็อกอิน", "บัญชี"},
			{"ปลอดภัย", "คุ้มครอง", "เชื่อถือ"},
			{"สนับสนุน", "ช่วยเหลือ", "บริการ"},
			{"ชำระ", "จ่าย", "เรียกเก็บ"},
			{"ธนาคาร", "การเงิน", "เงิน", "กระเป๋า"},
			{"จดหมาย", "อีเมล", "ข้อความ"},
			{"คลาวด์", "เซิร์ฟเวอร์", "โฮสติ้ง"},
			{"เว็บ", "ไซต์", "ออนไลน์", "เครือข่าย"},
			{"ข่าว", "บล็อก", "สื่อ"},
			{"ฟรี", "ของขวัญ", "โบนัส", "โปรโมชั่น"},
			{"ดีที่สุด", "ยอด", "พรีเมียม"},
			{"ใหม่", "ล่าสุด", "ทันสมัย"},
			{"เร็ว", "รวดเร็ว", "ด่วน"},
			{"แอป", "แอปพลิเคชัน", "ซอฟต์แวร์"},
			{"ดาวน์โหลด", "รับ", "ติดตั้ง"},
			{"ติดต่อ", "เชื่อมต่อ"},
			{"ค้นหา", "หา"},
			{"อัปเดต", "ต่ออายุ"},
		}
	case "vi":
		return [][]string{
			{"batdau", "khoidau"},
			{"cuahang", "shop", "thitruong", "mua"},
			{"dangnhap", "taikhoan"},
			{"antoan", "baove", "tincay"},
			{"hotro", "giupdo", "dichvu"},
			{"thanhtoan", "chitra", "hoadon"},
			{"nganhang", "taichinh", "tien", "vi"},
			{"thu", "email", "tinnhan"},
			{"dammay", "maychu", "hosting"},
			{"web", "trangweb", "online", "mang"},
			{"tintuc", "blog", "baochi", "truyenthong"},
			{"mienphi", "quatang", "thuong", "khuyenmai"},
			{"totnhat", "hangdau", "premium"},
			{"moi", "moinhat", "hiendai"},
			{"nhanh", "maule", "express"},
			{"ungdung", "phanmem", "congcu"},
			{"taive", "nhan", "caidat"},
			{"lienhe", "ketnoi"},
			{"timkiem", "tim"},
			{"capnhat", "giahan"},
		}
	case "la":
		return [][]string{
			{"initium", "principium", "exordium"},
			{"taberna", "mercatus", "emptio"},
			{"aditus", "ratio"},
			{"tutus", "securus", "fidelis"},
			{"auxilium", "adiutorium", "ministerium"},
			{"solutio", "pensio"},
			{"argentaria", "pecunia", "crumena"},
			{"epistula", "nuntius"},
			{"nubes", "servus", "hospitium"},
			{"tela", "situs", "rete", "porta"},
			{"nuntii", "acta"},
			{"gratis", "donum", "praemium"},
			{"optimus", "summus"},
			{"novus", "recens", "modernus"},
			{"celer", "velox"},
			{"instrumentum", "programma"},
			{"descendere", "accipere"},
			{"contactus", "coniungere"},
			{"quaerere", "invenire"},
		}
	case "hy":
		return [][]string{
			{"սկիզբ", "մեկնարկ"},
			{"խանութ", "շուկա", "գնում"},
			{"մուտք", "հաշիվ"},
			{"ապահով", "պաշտպանված", "վստահելի"},
			{"աջակցություն", "օգնություն", "ծառայություն"},
			{"վճարում", "վճար", "հաշիվ"},
			{"բանկ", "ֆինանս", "փող", "դրամապանակ"},
			{"փոստ", "նամակ", "հաղորդագրություն"},
			{"ամպ", "սերվեր", "հոսթինգ"},
			{"վեբ", "կայք", "առցանց", "ցանց"},
			{"նորություններ", "բլոգ", "մամուլ"},
			{"անվճար", "նվեր", "բոնուս"},
			{"լավագույն", "գագաթ", "պրեմիում"},
			{"նոր", "վերջին", "ժամանակակից"},
			{"արագ", "սրընթաց"},
			{"հավելված", "ծրագիր", "գործիք"},
			{"ներբեռնել", "ստանալ", "տեղադրել"},
			{"կապ", "միանալ"},
			{"որոնում", "գտնել"},
			{"թարմացնել", "նորացնել"},
		}
	case "ka":
		return [][]string{
			{"დაწყება", "დასაწყისი"},
			{"მაღაზია", "ბაზარი", "ყიდვა"},
			{"შესვლა", "ანგარიში"},
			{"უსაფრთხო", "დაცული", "სანდო"},
			{"მხარდაჭერა", "დახმარება", "სერვისი"},
			{"გადახდა", "ანგარიშსწორება"},
			{"ბანკი", "ფინანსები", "ფული", "საფულე"},
			{"ფოსტა", "წერილი", "შეტყობინება"},
			{"ღრუბელი", "სერვერი", "ჰოსტინგი"},
			{"ვები", "საიტი", "ონლაინ", "ქსელი"},
			{"სიახლეები", "ბლოგი", "პრესა", "მედია"},
			{"უფასო", "საჩუქარი", "ბონუსი"},
			{"საუკეთესო", "ტოპ", "პრემიუმ"},
			{"ახალი", "უახლესი", "თანამედროვე"},
			{"სწრაფი", "ჩქარი"},
			{"აპლიკაცია", "პროგრამა", "ხელსაწყო"},
			{"ჩამოტვირთვა", "მიღება", "დაყენება"},
			{"კონტაქტი", "დაკავშირება"},
			{"ძებნა", "პოვნა"},
		}
	case "ps":
		return [][]string{
			{"پیل", "شروع", "پیلامه"},
			{"پلورنځی", "بازار", "پیرودل"},
			{"ننوتل", "حساب"},
			{"خوندي", "ساتل", "باوري"},
			{"ملاتړ", "مرسته", "خدمت"},
			{"تادیه", "ورکړه", "بیلنګ"},
			{"بانک", "مالي", "پیسې", "کڅوړه"},
			{"پست", "بریښنالیک", "پیغام"},
			{"ورېځ", "سرور", "کوربتوب"},
			{"ویب", "پاڼه", "آنلاین", "شبکه"},
			{"خبرونه", "بلاګ", "مطبوعات", "رسنۍ"},
			{"وړیا", "ډالۍ", "انعام", "وړاندیز"},
			{"غوره", "لوړ", "ممتاز"},
			{"نوی", "وروستی", "عصري"},
			{"چټک", "ګړندی"},
			{"اپلیکیشن", "پروګرام", "وسیله"},
			{"ډاونلوډ", "ترلاسه", "نصب"},
			{"اړیکه", "وصل"},
			{"لټون", "موندل"},
			{"تازه", "نوی", "کول"},
		}
	default:
		// At least one safe line so the dataset isn't empty.
		return [][]string{{"start", "begin"}}
	}
}

// defaultPositive returns trust-signalling words used in combo squatting.
func defaultPositive(langID string) []string {
	switch langID {
	case "en":
		return []string{"good", "great", "excellent", "best", "safe", "trusted", "secure", "official", "verified", "fast", "easy", "free", "premium", "top", "new", "real", "genuine", "reliable", "quality", "happy"}
	case "fr":
		return []string{"bon", "super", "excellent", "meilleur", "sur", "fiable", "securise", "officiel", "verifie", "rapide", "facile", "gratuit", "premium", "top", "nouveau", "vrai", "authentique", "qualite", "heureux"}
	case "es":
		return []string{"bueno", "excelente", "mejor", "seguro", "fiable", "confiable", "oficial", "verificado", "rapido", "facil", "gratis", "premium", "top", "nuevo", "real", "autentico", "calidad", "feliz", "genial"}
	case "pt":
		return []string{"bom", "otimo", "excelente", "melhor", "seguro", "confiavel", "oficial", "verificado", "rapido", "facil", "gratis", "premium", "top", "novo", "real", "autentico", "qualidade", "feliz"}
	case "it":
		return []string{"buono", "ottimo", "eccellente", "migliore", "sicuro", "affidabile", "ufficiale", "verificato", "rapido", "facile", "gratis", "premium", "top", "nuovo", "vero", "autentico", "qualita", "felice"}
	case "de":
		return []string{"gut", "super", "sicher", "zuverlaessig", "toll", "beste", "offiziell", "verifiziert", "schnell", "einfach", "gratis", "premium", "top", "neu", "echt", "vertrauenswuerdig", "qualitaet", "gluecklich"}
	case "nl":
		return []string{"goed", "veilig", "betrouwbaar", "geweldig", "beste", "officieel", "geverifieerd", "snel", "makkelijk", "gratis", "premium", "top", "nieuw", "echt", "kwaliteit", "blij", "uitstekend"}
	case "sv":
		return []string{"bra", "saker", "palitlig", "utmarkt", "basta", "officiell", "verifierad", "snabb", "enkel", "gratis", "premium", "topp", "ny", "akta", "kvalitet", "glad"}
	case "no":
		return []string{"bra", "sikker", "palitelig", "utmerket", "beste", "offisiell", "verifisert", "rask", "enkel", "gratis", "premium", "topp", "ny", "ekte", "kvalitet", "glad"}
	case "da":
		return []string{"god", "sikker", "palidelig", "fremragende", "bedste", "officiel", "verificeret", "hurtig", "nem", "gratis", "premium", "top", "ny", "aegte", "kvalitet", "glad"}
	case "fi":
		return []string{"hyva", "turvallinen", "luotettava", "erinomainen", "paras", "virallinen", "vahvistettu", "nopea", "helppo", "ilmainen", "premium", "huippu", "uusi", "aito", "laatu", "iloinen"}
	case "ru":
		return []string{"хороший", "отличный", "лучший", "безопасный", "надежный", "официальный", "проверенный", "быстрый", "простой", "бесплатный", "премиум", "топ", "новый", "настоящий", "качество", "счастливый"}
	case "uk":
		return []string{"добрий", "відмінний", "кращий", "безпечний", "надійний", "офіційний", "перевірений", "швидкий", "простий", "безкоштовний", "преміум", "топ", "новий", "справжній", "якість", "щасливий"}
	case "pl":
		return []string{"dobry", "swietny", "najlepszy", "bezpieczny", "niezawodny", "oficjalny", "zweryfikowany", "szybki", "latwy", "darmowy", "premium", "top", "nowy", "prawdziwy", "jakosc", "szczesliwy"}
	case "cs":
		return []string{"dobry", "skvely", "nejlepsi", "bezpecny", "spolehlivy", "oficialni", "overeny", "rychly", "snadny", "zdarma", "premium", "top", "novy", "pravy", "kvalita", "stastny"}
	case "tr":
		return []string{"iyi", "harika", "en", "guvenli", "guvenilir", "resmi", "dogrulanmis", "hizli", "kolay", "ucretsiz", "premium", "top", "yeni", "gercek", "kalite", "mutlu"}
	case "el":
		return []string{"καλο", "τελειο", "καλυτερο", "ασφαλες", "αξιοπιστο", "επισημο", "επαληθευμενο", "γρηγορο", "ευκολο", "δωρεαν", "premium", "κορυφαιο", "νεο", "αληθινο", "ποιοτητα", "χαρουμενο"}
	case "ar":
		return []string{"جيد", "ممتاز", "افضل", "امن", "موثوق", "رسمي", "موثق", "سريع", "سهل", "مجاني", "مميز", "قمة", "جديد", "حقيقي", "اصلي", "جودة", "سعيد"}
	case "fa":
		return []string{"خوب", "عالی", "بهترین", "امن", "مطمئن", "رسمی", "تاییدشده", "سریع", "آسان", "رایگان", "ویژه", "برتر", "جدید", "واقعی", "اصل", "کیفیت", "خوشحال"}
	case "iw":
		return []string{"טוב", "מצוין", "הטוב", "בטוח", "אמין", "רשמי", "מאומת", "מהיר", "קל", "חינם", "פרימיום", "מוביל", "חדש", "אמיתי", "מקורי", "איכות", "שמח"}
	case "hi":
		return []string{"अच्छा", "उत्तम", "सर्वोत्तम", "सुरक्षित", "विश्वसनीय", "आधिकारिक", "सत्यापित", "तेज", "आसान", "मुफ्त", "प्रीमियम", "शीर्ष", "नया", "असली", "गुणवत्ता", "खुश"}
	case "zh":
		return []string{"好", "优秀", "最佳", "安全", "可靠", "官方", "已验证", "快速", "简单", "免费", "高级", "顶级", "新", "真实", "正品", "质量", "开心"}
	case "ja":
		return []string{"良い", "優秀", "最高", "安全", "信頼", "公式", "認証", "高速", "簡単", "無料", "プレミアム", "トップ", "新しい", "本物", "正規", "品質", "幸せ"}
	case "ko":
		return []string{"좋은", "훌륭한", "최고", "안전", "신뢰", "공식", "인증", "빠른", "쉬운", "무료", "프리미엄", "최상", "새로운", "진짜", "정품", "품질", "행복"}
	case "th":
		return []string{"ดี", "เยี่ยม", "ดีที่สุด", "ปลอดภัย", "น่าเชื่อถือ", "ทางการ", "ยืนยันแล้ว", "เร็ว", "ง่าย", "ฟรี", "พรีเมียม", "ยอด", "ใหม่", "จริง", "แท้", "คุณภาพ", "มีความสุข"}
	case "vi":
		return []string{"tot", "tuyet", "nhat", "antoan", "tincay", "chinhthuc", "daxacminh", "nhanh", "dedang", "mienphi", "premium", "hangdau", "moi", "that", "chinhhang", "chatluong", "vui"}
	case "la":
		return []string{"bonus", "optimus", "tutus", "fidelis", "officialis", "probatus", "celer", "facilis", "gratuitus", "summus", "novus", "verus", "germanus", "qualitas", "felix"}
	case "hy":
		return []string{"լավ", "գերազանց", "լավագույն", "ապահով", "վստահելի", "պաշտոնական", "ստուգված", "արագ", "հեշտ", "անվճար", "պրեմիում", "նոր", "իսկական", "որակ", "ուրախ"}
	case "ka":
		return []string{"კარგი", "შესანიშნავი", "საუკეთესო", "უსაფრთხო", "სანდო", "ოფიციალური", "დადასტურებული", "სწრაფი", "მარტივი", "უფასო", "პრემიუმ", "ახალი", "ნამდვილი", "ხარისხი", "ბედნიერი"}
	case "ps":
		return []string{"ښه", "عالي", "غوره", "خوندي", "باوري", "رسمي", "تایید", "شوی", "چټک", "اسان", "وړیا", "ممتاز", "نوی", "ریښتینی", "کیفیت", "خوشحاله"}
	default:
		return []string{"good", "safe", "trusted", "official", "verified"}
	}
}

// defaultNegative returns alarm words used in combo squatting and in scoring.
func defaultNegative(langID string) []string {
	switch langID {
	case "en":
		return []string{"bad", "unsafe", "fake", "scam", "malicious", "fraud", "phishing", "spam", "virus", "malware", "danger", "warning", "error", "fail", "blocked", "hacked", "stolen", "illegal", "suspicious", "untrusted"}
	case "fr":
		return []string{"mauvais", "dangereux", "faux", "arnaque", "malveillant", "fraude", "hameconnage", "spam", "virus", "logiciel", "danger", "avertissement", "erreur", "echec", "bloque", "pirate", "vole", "illegal", "suspect"}
	case "es":
		return []string{"malo", "inseguro", "falso", "estafa", "malicioso", "fraude", "phishing", "spam", "virus", "malware", "peligro", "advertencia", "error", "fallo", "bloqueado", "hackeado", "robado", "ilegal", "sospechoso"}
	case "pt":
		return []string{"mau", "inseguro", "falso", "golpe", "malicioso", "fraude", "phishing", "spam", "virus", "malware", "perigo", "aviso", "erro", "falha", "bloqueado", "hackeado", "roubado", "ilegal", "suspeito"}
	case "it":
		return []string{"cattivo", "insicuro", "falso", "truffa", "malevolo", "frode", "phishing", "spam", "virus", "malware", "pericolo", "avviso", "errore", "fallito", "bloccato", "violato", "rubato", "illegale", "sospetto"}
	case "de":
		return []string{"schlecht", "unsicher", "gefaelscht", "betrug", "boesartig", "phishing", "spam", "virus", "schadsoftware", "gefahr", "warnung", "fehler", "fehlschlag", "gesperrt", "gehackt", "gestohlen", "illegal", "verdaechtig"}
	case "nl":
		return []string{"slecht", "onveilig", "nep", "oplichting", "kwaadaardig", "fraude", "phishing", "spam", "virus", "malware", "gevaar", "waarschuwing", "fout", "mislukt", "geblokkeerd", "gehackt", "gestolen", "illegaal", "verdacht"}
	case "sv":
		return []string{"dalig", "osaker", "falsk", "bedrageri", "skadlig", "phishing", "spam", "virus", "malware", "fara", "varning", "fel", "misslyckad", "blockerad", "hackad", "stulen", "olaglig", "misstankt"}
	case "no":
		return []string{"darlig", "usikker", "falsk", "svindel", "skadelig", "phishing", "spam", "virus", "malware", "fare", "advarsel", "feil", "mislykket", "blokkert", "hacket", "stjalet", "ulovlig", "mistenkelig"}
	case "da":
		return []string{"darlig", "usikker", "falsk", "svindel", "ondsindet", "phishing", "spam", "virus", "malware", "fare", "advarsel", "fejl", "mislykket", "blokeret", "hacket", "stjalet", "ulovlig", "mistaenkelig"}
	case "fi":
		return []string{"huono", "turvaton", "vaarennos", "huijaus", "haitallinen", "phishing", "roskaposti", "virus", "haittaohjelma", "vaara", "varoitus", "virhe", "epaonnistui", "estetty", "murrettu", "varastettu", "laiton", "epailyttava"}
	case "ru":
		return []string{"плохой", "небезопасный", "поддельный", "мошенничество", "вредоносный", "фишинг", "спам", "вирус", "вредонос", "опасность", "предупреждение", "ошибка", "сбой", "заблокирован", "взломан", "украден", "незаконный", "подозрительный"}
	case "uk":
		return []string{"поганий", "небезпечний", "підроблений", "шахрайство", "шкідливий", "фішинг", "спам", "вірус", "загроза", "попередження", "помилка", "збій", "заблокований", "зламаний", "вкрадений", "незаконний", "підозрілий"}
	case "pl":
		return []string{"zly", "niebezpieczny", "falszywy", "oszustwo", "zlosliwy", "phishing", "spam", "wirus", "malware", "zagrozenie", "ostrzezenie", "blad", "awaria", "zablokowany", "zhakowany", "skradziony", "nielegalny", "podejrzany"}
	case "cs":
		return []string{"spatny", "nebezpecny", "falesny", "podvod", "skodlivy", "phishing", "spam", "virus", "malware", "nebezpeci", "varovani", "chyba", "selhani", "blokovano", "hacknuto", "ukradeno", "nelegalni", "podezrely"}
	case "tr":
		return []string{"kotu", "guvensiz", "sahte", "dolandiricilik", "zararli", "oltalama", "spam", "virus", "zararliyazilim", "tehlike", "uyari", "hata", "basarisiz", "engellendi", "hacklendi", "calindi", "yasadisi", "supheli"}
	case "el":
		return []string{"κακο", "επισφαλες", "ψευτικο", "απατη", "κακοβουλο", "ηλεκτρονικο", "spam", "ιος", "κακολογισμικο", "κινδυνος", "προειδοποιηση", "σφαλμα", "αποτυχια", "αποκλεισμενο", "παραβιασμενο", "κλεμμενο", "παρανομο", "υποπτο"}
	case "ar":
		return []string{"سيء", "غيرآمن", "مزيف", "احتيال", "ضار", "تصيد", "بريدمزعج", "فيروس", "برمجيات", "خطر", "تحذير", "خطأ", "فشل", "محظور", "اختراق", "مسروق", "غيرقانوني", "مشبوه"}
	case "fa":
		return []string{"بد", "ناامن", "جعلی", "کلاهبرداری", "مخرب", "فیشینگ", "هرزنامه", "ویروس", "بدافزار", "خطر", "هشدار", "خطا", "شکست", "مسدود", "هک", "دزدیده", "غیرقانونی", "مشکوک"}
	case "iw":
		return []string{"רע", "לאבטוח", "מזויף", "הונאה", "זדוני", "דיוג", "ספאם", "וירוס", "נוזקה", "סכנה", "אזהרה", "שגיאה", "כשל", "חסום", "נפרץ", "גנוב", "לאחוקי", "חשוד"}
	case "hi":
		return []string{"बुरा", "असुरक्षित", "नकली", "धोखा", "दुर्भावनापूर्ण", "फ़िशिंग", "स्पैम", "वायरस", "मैलवेयर", "खतरा", "चेतावनी", "त्रुटि", "विफल", "अवरुद्ध", "हैक", "चोरी", "अवैध", "संदिग्ध"}
	case "zh":
		return []string{"坏", "不安全", "假", "诈骗", "恶意", "钓鱼", "垃圾", "病毒", "恶意软件", "危险", "警告", "错误", "失败", "封锁", "被黑", "被盗", "非法", "可疑"}
	case "ja":
		return []string{"悪い", "危険", "偽物", "詐欺", "悪意", "フィッシング", "スパム", "ウイルス", "マルウェア", "警告", "エラー", "失敗", "ブロック", "ハッキング", "盗難", "違法", "不審"}
	case "ko":
		return []string{"나쁜", "위험", "가짜", "사기", "악성", "피싱", "스팸", "바이러스", "멀웨어", "경고", "오류", "실패", "차단", "해킹", "도난", "불법", "의심"}
	case "th":
		return []string{"แย่", "ไม่ปลอดภัย", "ปลอม", "หลอกลวง", "ประสงค์ร้าย", "ฟิชชิ่ง", "สแปม", "ไวรัส", "มัลแวร์", "อันตราย", "คำเตือน", "ข้อผิดพลาด", "ล้มเหลว", "ถูกบล็อก", "ถูกแฮก", "ถูกขโมย", "ผิดกฎหมาย", "น่าสงสัย"}
	case "vi":
		return []string{"xau", "khongantoan", "gia", "luadao", "doche", "phishing", "spam", "virus", "malware", "nguyhiem", "canhbao", "loi", "thatbai", "bichan", "bihack", "bidanhcap", "batgiaphap", "dangngo"}
	case "la":
		return []string{"malus", "periculosus", "falsus", "fraus", "malignus", "virus", "periculum", "monitum", "error", "defectus", "interclusus", "furatus", "illicitus", "suspectus"}
	case "hy":
		return []string{"վատ", "վտանգավոր", "կեղծ", "խարդախություն", "վնասակար", "ֆիշինգ", "սպամ", "վիրուս", "վտանգ", "զգուշացում", "սխալ", "ձախողում", "արգելափակված", "կոտրված", "գողացված", "անօրինական", "կասկածելի"}
	case "ka":
		return []string{"ცუდი", "საშიში", "ყალბი", "თაღლითობა", "მავნე", "ფიშინგი", "სპამი", "ვირუსი", "საფრთხე", "გაფრთხილება", "შეცდომა", "წარუმატებელი", "დაბლოკილი", "გატეხილი", "მოპარული", "უკანონო", "საეჭვო"}
	case "ps":
		return []string{"بد", "ناخوندي", "جعلي", "درغلي", "زیانمن", "فیشینګ", "سپیم", "ویروس", "ګواښ", "خبرداری", "تېروتنه", "ناکامي", "بند", "شوی", "هیک", "غلا", "غیرقانوني", "شکمن"}
	default:
		return []string{"bad", "unsafe", "fake", "scam", "malicious"}
	}
}

// defaultStopwords returns high-frequency function words, used to strip filler
// from tokenized names before generating variants.
func defaultStopwords(langID string) []string {
	switch langID {
	case "en":
		return []string{"a", "an", "the", "and", "or", "but", "if", "while", "of", "to", "in", "on", "at", "by", "for", "with", "about", "from", "up", "down", "out", "over", "under", "again", "then", "once", "here", "there", "when", "where", "why", "how", "all", "any", "both", "each", "few", "more", "most", "other", "some", "such", "no", "nor", "not", "only", "own", "same", "so", "than", "too", "very", "can", "will", "just", "is", "are", "was", "were", "be", "been"}
	case "fr":
		return []string{"le", "la", "les", "un", "une", "des", "de", "du", "au", "aux", "et", "ou", "mais", "donc", "or", "ni", "car", "ce", "cet", "cette", "ces", "mon", "ton", "son", "notre", "votre", "leur", "je", "tu", "il", "elle", "nous", "vous", "ils", "elles", "que", "qui", "quoi", "dont", "pour", "par", "avec", "sans", "sous", "sur", "dans", "chez", "vers", "entre", "est", "sont", "etait", "sera", "ne", "pas", "plus", "tres", "tout"}
	case "es":
		return []string{"el", "la", "los", "las", "un", "una", "unos", "unas", "de", "del", "al", "y", "o", "pero", "sino", "porque", "que", "quien", "cual", "cuyo", "mi", "tu", "su", "nuestro", "vuestro", "yo", "ella", "nosotros", "vosotros", "ellos", "para", "por", "con", "sin", "sobre", "bajo", "entre", "hacia", "desde", "hasta", "es", "son", "era", "sera", "no", "mas", "muy", "todo"}
	case "pt":
		return []string{"o", "a", "os", "as", "um", "uma", "uns", "umas", "de", "do", "da", "dos", "das", "no", "na", "nos", "nas", "e", "ou", "mas", "porque", "que", "quem", "qual", "cujo", "meu", "teu", "seu", "nosso", "vosso", "eu", "ele", "ela", "vos", "eles", "elas", "para", "por", "com", "sem", "sobre", "sob", "entre", "ate", "desde", "sao", "era", "sera", "nao", "mais", "muito", "todo"}
	case "it":
		return []string{"il", "lo", "la", "i", "gli", "le", "un", "uno", "una", "di", "del", "della", "dei", "delle", "e", "o", "ma", "perche", "che", "chi", "quale", "cui", "mio", "tuo", "suo", "nostro", "vostro", "loro", "io", "tu", "egli", "noi", "voi", "essi", "per", "con", "senza", "su", "sotto", "tra", "fra", "da", "a", "in", "sono", "era", "sara", "non", "piu", "molto", "tutto"}
	case "de":
		return []string{"der", "die", "das", "ein", "eine", "einen", "einem", "eines", "und", "oder", "aber", "denn", "weil", "dass", "wer", "was", "welche", "mein", "dein", "sein", "unser", "euer", "ihr", "ich", "du", "er", "sie", "es", "wir", "fuer", "mit", "ohne", "auf", "unter", "zwischen", "von", "zu", "in", "an", "bei", "ist", "sind", "war", "wird", "nicht", "mehr", "sehr", "alle"}
	case "nl":
		return []string{"de", "het", "een", "en", "of", "maar", "want", "omdat", "dat", "wie", "wat", "welke", "mijn", "jouw", "zijn", "haar", "onze", "jullie", "ik", "jij", "hij", "zij", "wij", "voor", "met", "zonder", "op", "onder", "tussen", "van", "naar", "in", "aan", "bij", "is", "was", "wordt", "niet", "meer", "zeer", "alle"}
	case "sv":
		return []string{"en", "ett", "den", "det", "de", "och", "eller", "men", "for", "att", "som", "vem", "vad", "vilken", "min", "din", "hans", "hennes", "var", "era", "jag", "du", "han", "hon", "vi", "ni", "med", "utan", "pa", "under", "mellan", "fran", "till", "i", "av", "vid", "ar", "blir", "inte", "mer", "mycket", "alla"}
	case "no":
		return []string{"en", "et", "den", "det", "de", "og", "eller", "men", "for", "at", "som", "hvem", "hva", "hvilken", "min", "din", "hans", "hennes", "var", "deres", "jeg", "du", "han", "hun", "vi", "dere", "med", "uten", "pa", "under", "mellom", "fra", "til", "i", "av", "ved", "er", "blir", "ikke", "mer", "veldig", "alle"}
	case "da":
		return []string{"en", "et", "den", "det", "de", "og", "eller", "men", "for", "at", "som", "hvem", "hvad", "hvilken", "min", "din", "hans", "hendes", "vores", "jeres", "jeg", "du", "han", "hun", "vi", "med", "uden", "pa", "under", "mellem", "fra", "til", "i", "af", "ved", "er", "var", "bliver", "ikke", "mere", "meget", "alle"}
	case "fi":
		return []string{"ja", "tai", "mutta", "silla", "etta", "joka", "mika", "kuka", "minun", "sinun", "hanen", "meidan", "teidan", "heidan", "mina", "sina", "han", "me", "te", "he", "kanssa", "ilman", "paalla", "alla", "valissa", "sta", "on", "oli", "tulee", "ei", "enemman", "hyvin", "kaikki"}
	case "ru":
		return []string{"и", "или", "но", "что", "чтобы", "как", "кто", "где", "когда", "мой", "твой", "его", "её", "наш", "ваш", "их", "я", "ты", "он", "она", "мы", "вы", "они", "для", "с", "без", "на", "под", "между", "от", "до", "в", "к", "при", "это", "тот", "этот", "не", "более", "очень", "все"}
	case "uk":
		return []string{"і", "або", "але", "що", "щоб", "як", "хто", "де", "коли", "мій", "твій", "його", "її", "наш", "ваш", "їх", "я", "ти", "він", "вона", "ми", "ви", "вони", "для", "з", "без", "на", "під", "між", "від", "до", "в", "к", "при", "це", "той", "цей", "не", "більш", "дуже", "всі"}
	case "pl":
		return []string{"i", "lub", "ale", "ze", "aby", "jak", "kto", "gdzie", "kiedy", "moj", "twoj", "jego", "jej", "nasz", "wasz", "ich", "ja", "ty", "on", "ona", "my", "wy", "oni", "dla", "z", "bez", "na", "pod", "miedzy", "od", "do", "w", "przy", "to", "ten", "nie", "bardziej", "bardzo", "wszystko"}
	case "cs":
		return []string{"a", "nebo", "ale", "ze", "aby", "jak", "kdo", "kde", "kdy", "muj", "tvuj", "jeho", "jeji", "nas", "vas", "jejich", "ja", "ty", "on", "ona", "my", "vy", "oni", "pro", "s", "bez", "na", "pod", "mezi", "od", "do", "v", "pri", "to", "ten", "ne", "vice", "velmi", "vse"}
	case "tr":
		return []string{"ve", "veya", "ama", "ki", "icin", "gibi", "kim", "nerede", "ne", "zaman", "benim", "senin", "onun", "bizim", "sizin", "onlarin", "ben", "sen", "o", "biz", "siz", "onlar", "ile", "olmadan", "uzerinde", "altinda", "arasinda", "dan", "den", "bir", "bu", "su", "degil", "daha", "cok", "hepsi"}
	case "el":
		return []string{"και", "η", "το", "ο", "οι", "τα", "ενα", "μια", "αλλα", "οτι", "για", "οπως", "ποιος", "που", "ποτε", "μου", "σου", "του", "της", "μας", "σας", "τους", "εγω", "εσυ", "αυτος", "αυτη", "εμεις", "εσεις", "αυτοι", "με", "χωρις", "πανω", "κατω", "μεταξυ", "απο", "σε", "ειναι", "ηταν", "δεν", "πιο", "πολυ", "ολα"}
	case "ar":
		return []string{"و", "او", "لكن", "ان", "كي", "مثل", "من", "اين", "متى", "لي", "لك", "له", "لها", "لنا", "لكم", "لهم", "انا", "انت", "هو", "هي", "نحن", "انتم", "هم", "مع", "بدون", "على", "تحت", "بين", "الى", "في", "عند", "هذا", "ذلك", "ليس", "اكثر", "جدا", "كل"}
	case "fa":
		return []string{"و", "یا", "اما", "که", "تا", "مانند", "از", "کجا", "کی", "من", "تو", "او", "ما", "شما", "ایشان", "با", "بدون", "روی", "زیر", "بین", "به", "در", "نزد", "این", "آن", "نیست", "بیشتر", "خیلی", "همه"}
	case "iw":
		return []string{"ו", "או", "אבל", "כי", "כמו", "מי", "איפה", "מתי", "שלי", "שלך", "שלו", "שלה", "שלנו", "שלכם", "שלהם", "אני", "אתה", "הוא", "היא", "אנחנו", "אתם", "הם", "עם", "בלי", "על", "תחת", "בין", "מ", "ל", "ב", "אצל", "זה", "לא", "יותר", "מאוד", "כל"}
	case "hi":
		return []string{"और", "या", "लेकिन", "कि", "जैसे", "कौन", "कहाँ", "कब", "मेरा", "तुम्हारा", "उसका", "हमारा", "आपका", "उनका", "मैं", "तुम", "वह", "हम", "आप", "वे", "के", "साथ", "बिना", "पर", "नीचे", "बीच", "से", "को", "में", "पास", "यह", "नहीं", "अधिक", "बहुत", "सब"}
	case "zh":
		return []string{"的", "了", "和", "或", "但是", "因为", "如果", "谁", "哪里", "什么", "我", "你", "他", "她", "我们", "你们", "他们", "与", "没有", "在", "上", "下", "之间", "从", "到", "这", "那", "不", "更", "很", "都"}
	case "ja":
		return []string{"の", "に", "は", "を", "が", "と", "も", "で", "から", "まで", "より", "へ", "や", "か", "ね", "よ", "です", "ます", "これ", "それ", "あれ", "この", "その", "あの", "わたし", "あなた", "かれ", "かのじょ", "ない", "もっと", "とても", "すべて"}
	case "ko":
		return []string{"의", "에", "는", "을", "가", "와", "도", "로", "부터", "까지", "보다", "이", "그", "저", "나", "너", "그녀", "우리", "너희", "그들", "없이", "위", "아래", "사이", "아니다", "더", "매우", "모두"}
	case "th":
		return []string{"และ", "หรือ", "แต่", "ที่", "เพราะ", "ถ้า", "ใคร", "ที่ไหน", "อะไร", "ฉัน", "คุณ", "เขา", "เธอ", "เรา", "พวกเขา", "กับ", "ไม่มี", "บน", "ใต้", "ระหว่าง", "จาก", "ถึง", "ใน", "นี้", "นั้น", "ไม่", "มากกว่า", "มาก", "ทั้งหมด"}
	case "vi":
		return []string{"va", "hoac", "nhung", "ma", "vi", "neu", "ai", "dau", "gi", "toi", "ban", "anh", "chi", "chung", "ho", "voi", "khong", "tren", "duoi", "giua", "tu", "den", "trong", "nay", "do", "hon", "rat", "tatca"}
	case "la":
		return []string{"et", "aut", "sed", "quia", "si", "qui", "quae", "quod", "ubi", "quando", "meus", "tuus", "suus", "noster", "vester", "ego", "tu", "ille", "illa", "nos", "vos", "illi", "cum", "sine", "super", "sub", "inter", "ab", "ad", "in", "apud", "hic", "non", "magis", "valde", "omnis"}
	case "hy":
		return []string{"և", "կամ", "բայց", "որ", "եթե", "ով", "որտեղ", "երբ", "իմ", "քո", "նրա", "մեր", "ձեր", "նրանց", "ես", "դու", "նա", "մենք", "դուք", "նրանք", "հետ", "առանց", "վրա", "տակ", "միջև", "ից", "ին", "մեջ", "այս", "այն", "ոչ", "ավելի", "շատ", "բոլոր"}
	case "ka":
		return []string{"და", "ან", "მაგრამ", "რომ", "თუ", "ვინ", "სად", "როდის", "ჩემი", "შენი", "მისი", "ჩვენი", "თქვენი", "მათი", "მე", "შენ", "ის", "ჩვენ", "თქვენ", "ისინი", "თან", "გარეშე", "ზე", "ქვეშ", "შორის", "დან", "ში", "ეს", "არა", "მეტი", "ძალიან", "ყველა"}
	case "ps":
		return []string{"او", "یا", "خو", "چې", "که", "څوک", "چیرته", "کله", "زما", "ستا", "د", "هغه", "زموږ", "ستاسو", "هغوی", "زه", "ته", "موږ", "تاسو", "سره", "پرته", "په", "لاندې", "منځ", "له", "کې", "دا", "نه", "ډیر", "ټول"}
	default:
		return []string{"the", "and", "or", "of", "to", "in", "for", "with"}
	}
}

var SyncLanguagesCmd = cli.Command{
	Name:  "sync-languages",
	Usage: "Generate datasets/languages/<lang>/ files from registered language plugins",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "dir",
			Value: "datasets/languages",
			Usage: "output directory for language datasets",
		},
		&cli.BoolFlag{
			Name:  "overwrite",
			Value: false,
			Usage: "overwrite existing dataset files",
		},
	},
	Action: func(cCtx *cli.Context) error {
		base := cCtx.String("dir")
		overwrite := cCtx.Bool("overwrite")

		langs := languages.Languages()
		sort.Slice(langs, func(i, j int) bool { return langs[i].Id() < langs[j].Id() })

		for _, l := range langs {
			langDir := filepath.Join(base, l.Id())
			if err := os.MkdirAll(langDir, 0o755); err != nil {
				return err
			}

			// Write helpers
			write := func(filename string, content []byte) error {
				path := filepath.Join(langDir, filename)
				if !overwrite {
					if st, err := os.Stat(path); err == nil {
						// Do not overwrite curated datasets, but allow filling files that are empty (size=0).
						if st.Size() > 0 {
							return nil
						}
					}
				}
				return os.WriteFile(path, content, 0o644)
			}
			writeText := func(filename, content string) error {
				if content != "" && !strings.HasSuffix(content, "\n") {
					content += "\n"
				}
				return write(filename, []byte(content))
			}

			// Aggregate tokens that "belong" to the language plugin datasets to build a minimal word.lst.
			words := map[string]bool{}
			add := func(s string) {
				s = strings.TrimSpace(s)
				if s == "" {
					return
				}
				words[s] = true
			}

			// Numerals: "number token token ..."
			{
				var keys []string
				for k := range l.Numerals() {
					keys = append(keys, k)
				}
				sort.Slice(keys, func(i, j int) bool {
					ai, err1 := strconv.Atoi(keys[i])
					aj, err2 := strconv.Atoi(keys[j])
					if err1 == nil && err2 == nil {
						return ai < aj
					}
					return keys[i] < keys[j]
				})

				var b strings.Builder
				for _, k := range keys {
					b.WriteString(k)
					add(k)
					for _, tok := range l.Numerals()[k] {
						if strings.TrimSpace(tok) == "" {
							continue
						}
						b.WriteString(" ")
						b.WriteString(tok)
						add(tok)
					}
					b.WriteString("\n")
				}
				if err := writeText("numeral.lst", b.String()); err != nil {
					return err
				}
			}

			// Graphemes/vowels: one per line
			{
				var b strings.Builder
				for _, g := range l.Graphemes() {
					g = strings.TrimSpace(g)
					if g == "" {
						continue
					}
					b.WriteString(g)
					b.WriteString("\n")
					add(g)
				}
				if err := writeText("grapheme.lst", b.String()); err != nil {
					return err
				}
			}
			{
				var b strings.Builder
				for _, v := range l.Vowels() {
					v = strings.TrimSpace(v)
					if v == "" {
						continue
					}
					b.WriteString(v)
					b.WriteString("\n")
					add(v)
				}
				if err := writeText("vowel.lst", b.String()); err != nil {
					return err
				}
			}

			// Misspellings / homophones: one set per line
			{
				var b strings.Builder
				for _, set := range l.Misspellings() {
					var parts []string
					for _, tok := range set {
						tok = strings.TrimSpace(tok)
						if tok != "" {
							parts = append(parts, tok)
							add(tok)
						}
					}
					if len(parts) == 0 {
						continue
					}
					b.WriteString(strings.Join(parts, " "))
					b.WriteString("\n")
				}
				if err := writeText("misspelling.lst", b.String()); err != nil {
					return err
				}
			}
			{
				var b strings.Builder
				for _, set := range l.Homophones() {
					var parts []string
					for _, tok := range set {
						tok = strings.TrimSpace(tok)
						if tok != "" {
							parts = append(parts, tok)
							add(tok)
						}
					}
					if len(parts) == 0 {
						continue
					}
					b.WriteString(strings.Join(parts, " "))
					b.WriteString("\n")
				}
				if err := writeText("homophone.lst", b.String()); err != nil {
					return err
				}
			}

			// Antonyms: "word antonym1 antonym2 ..."
			{
				var keys []string
				for k := range l.Antonyms() {
					keys = append(keys, k)
				}
				sort.Strings(keys)

				var b strings.Builder
				for _, k := range keys {
					k = strings.TrimSpace(k)
					if k == "" {
						continue
					}
					add(k)
					b.WriteString(k)
					for _, a := range l.Antonyms()[k] {
						a = strings.TrimSpace(a)
						if a == "" {
							continue
						}
						add(a)
						b.WriteString(" ")
						b.WriteString(a)
					}
					b.WriteString("\n")
				}
				if err := writeText("antonym.lst", b.String()); err != nil {
					return err
				}
			}

			// Homoglyphs: "char homoglyph1 homoglyph2 ..."
			{
				var keys []string
				for k := range l.Homoglyphs() {
					keys = append(keys, k)
				}
				sort.Strings(keys)

				var buf bytes.Buffer
				for _, k := range keys {
					k = strings.TrimSpace(k)
					if k == "" {
						continue
					}
					add(k)
					buf.WriteString(k)
					for _, h := range l.Homoglyphs()[k] {
						h = strings.TrimSpace(h)
						if h == "" {
							continue
						}
						add(h)
						buf.WriteString(" ")
						buf.WriteString(h)
					}
					buf.WriteString("\n")
				}
				if err := write("homoglyph.lst", buf.Bytes()); err != nil {
					return err
				}
			}

			// Optional / currently unused by the importer, but keep directory parity with existing datasets.
			{
				var all []string
				for w := range words {
					all = append(all, w)
				}
				sort.Strings(all)
				var b strings.Builder
				for _, w := range all {
					b.WriteString(w)
					b.WriteString("\n")
				}
				if err := writeText("word.lst", b.String()); err != nil {
					return err
				}
			}
			{
				// One word per line
				var b strings.Builder
				for _, w := range defaultStopwords(l.Id()) {
					w = strings.TrimSpace(w)
					if w == "" {
						continue
					}
					add(w)
					b.WriteString(w)
					b.WriteString("\n")
				}
				if err := writeText("stopword.lst", b.String()); err != nil {
					return err
				}
			}
			{
				// One word per line
				var b strings.Builder
				for _, w := range defaultPositive(l.Id()) {
					w = strings.TrimSpace(w)
					if w == "" {
						continue
					}
					add(w)
					b.WriteString(w)
					b.WriteString("\n")
				}
				if err := writeText("positive.lst", b.String()); err != nil {
					return err
				}
			}
			{
				// One word per line
				var b strings.Builder
				for _, w := range defaultNegative(l.Id()) {
					w = strings.TrimSpace(w)
					if w == "" {
						continue
					}
					add(w)
					b.WriteString(w)
					b.WriteString("\n")
				}
				if err := writeText("negative.lst", b.String()); err != nil {
					return err
				}
			}
			{
				var b strings.Builder
				for _, set := range defaultSynonyms(l.Id()) {
					var parts []string
					for _, tok := range set {
						tok = strings.TrimSpace(tok)
						if tok != "" {
							parts = append(parts, tok)
							add(tok)
						}
					}
					if len(parts) == 0 {
						continue
					}
					b.WriteString(strings.Join(parts, " "))
					b.WriteString("\n")
				}
				if err := writeText("synonym.lst", b.String()); err != nil {
					return err
				}
			}
			if err := writeText("token.lst", ""); err != nil {
				return err
			}
		}

		fmt.Printf("Synced %d language dataset folders into %s (overwrite=%v)\n", len(langs), base, overwrite)
		return nil
	},
}
