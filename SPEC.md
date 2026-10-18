# Disiplin.exe — Legacy .NET Application Technical Specification (`SPEC.md`)

## 1. Overview & Binary Inventory

**Disiplin** (`Disiplin.exe`, v1.0.0.0, copyright *YAZILIM EVİ - 2016*, updated through 2025/2026 MEB regulation amendments) is a 32-bit Windows Forms utility written in VB.NET (compiled to .NET Framework 4.0) for Turkish high schools (**LİSE** — *Ödül ve Disiplin Kurulu*) and middle schools (**ORTAOKUL** — *Öğrenci Davranışlarını Değerlendirme Kurulu*).

### 1.1 Unpacked Files in `./disiplin/`
| File | Type | Role |
|---|---|---|
| `Disiplin.exe` | .NET 4.0 WinForms (VB.NET) | Main application executable containing 27 forms/modules and embedded resources. Decompiled to `./extracted/src/Disiplin/`. |
| `Disiplin.pdb` | PDB Symbols | Original debugging symbols from `E:\PROGRAMLAR\Disiplin - 2025\Disiplin\obj\x86\Debug\Disiplin.pdb`. |
| `Disiplin.xml` | XML Doc | VB.NET compiler XML documentation file. |
| `PswNET.dll` | .NET 2.0 Assembly | Managed wrapper (`PswNET.MyTools.PassWord`) that P/Invokes `PswD.dll` and computes hardware-bound activation codes. Decompiled to `./extracted/src/PswNET/`. |
| `PswD.dll` | Native 32-bit Win32 PE DLL | Native DLL exporting `PassWord` (`@PassWord$qqsrx17System@AnsiStringt1`). Used only by `PswNET.dll` for legacy license verification. |
| `dsp.mdb` | Microsoft Jet 4.0 (Access 2000) | Password-protected local database storing all school settings, regulation clauses, students, incidents, verdicts, agenda items, and meeting minutes. Extracted to `./extracted/db/`. |

### 1.2 License / Demo Mode Mechanism (`Module_Pasword.cs`, `Frm17Hakkinda.cs`, `PswNET.dll`)
- `Module_Pasword.DemoKontrol()` reads `KurumKodu` and `OkulAdı` from `Genelbilgiler` and checks the Windows Registry key `HKEY_CURRENT_USER\Software\VB and VBA Program Settings\DİSİPLİN\GÜVENLİK\ŞİFRE`.
- It computes `PassWord.PassWord("DSP", KurumKodu + " " + OkulAdı)`. If the registry value matches, `Module_Pasword.DemoBoolean = true` (licensed mode).
- When `DemoBoolean == false` (unlicensed/demo mode):
  - Certain print actions open `Frm17Hakkinda` (license activation dialog) and abort, OR generate the PDF and stamp `img.jpg` centered underneath every page into `DİSİPLİNdemorapor.pdf` while deleting the original PDF.

---

## 2. Database Architecture (`dsp.mdb`)

### 2.1 Connection & Password
- **Engine**: Microsoft Jet 4.0 (`Provider=Microsoft.Jet.OLEDB.4.0`)
- **File Path**: `{Application.StartupPath}\dsp.mdb`
- **Database Password**: `17100301`
- **Connection String** (`dsp_Modül.Baglan()`):
  ```text
  Provider=Microsoft.Jet.OLEDB.4.0;Data Source=dsp.mdb;Jet OLEDB:Database Password=17100301;
  ```

### 2.2 Runtime Schema Migration (`dsp_Modül.TablolarınKontrolü()`)
On startup and after restore/school-type changes, `TablolarınKontrolü()` checks `ADODB.Connection.OpenSchema(adSchemaTables)` and `OpenSchema(adSchemaColumns)` and automatically creates missing tables or `ALTER TABLE ... ADD COLUMN ... TEXT(255)` for any missing columns. All columns across all user tables are stored as Jet `TEXT` (Unicode variable-length text, max 255 chars, except memo/text fields created without length limit).

### 2.3 Complete Table Schemas

#### 1. `Genelbilgiler` (1 row — School & Board Configuration)
| Column | Description | Default / Example |
|---|---|---|
| `sıra` | Row ID | `null` |
| `okultürü` | School type (`"LİSE"` or `"ORTAOKUL"`) | `"LİSE"` |
| `ÖgretimYılı` | Active academic year (`"YYYY - YYYY"` or `"YYYY-YYYY"`) | `"2026-2027"` |
| `kurumkodu` | MEB Institution Code | e.g. `"123456"` |
| `il` | Province (Uppercase) | e.g. `"ANKARA"` |
| `ilçe` | District (Uppercase) | e.g. `"ÇANKAYA"` |
| `OkulAdı` | School Name | e.g. `"Atatürk Anadolu Lisesi"` |
| `Müdüradı` | School Principal Name | |
| `Antet` | Multi-line official letterhead (`T.C.\n...`) | |
| `ResmiYazıKodu` | Official correspondence code (`Sayı`) | |
| `adres` | School street address | |
| `telno` | School phone number | |
| `faxno` | School fax number | |
| `uyarisayi` | Last sequential number used for Written Warnings (`Frm06YaziliUyari`) | `"1"` |
| `KurulBaşkanı` | Primary Board Chair name | |
| `KurulBaşkanıünvan` | Primary Board Chair title (e.g. `"Müdür Yardımcısı"`) | |
| `Üye1` .. `Üye4` | Primary Board Members 1–4 names | |
| `Üye1ünvan` .. `Üye4ünvan` | Primary Board Members 1–4 titles/branches | |
| `YedekKurulBaşkanı` | Substitute Board Chair name | |
| `YedekKurulBaşkanıünvan` | Substitute Board Chair title | |
| `YedekÜye1` .. `YedekÜye4` | Substitute Board Members 1–4 names | |
| `YedekÜye1ünvan` .. `YedekÜye4ünvan` | Substitute Board Members 1–4 titles/branches | |

*Derived Runtime Globals in `dsp_Modül`*:
- When `okultürü == "LİSE"` -> `dsp_Modül.KurulAdı = "Ödül ve Disiplin"`
- When `okultürü == "ORTAOKUL"` -> `dsp_Modül.KurulAdı = "Öğrenci Davranışlarını Değerlendirme"`

#### 2. `cezamaddeleri` (119 rows — Regulation Clauses)
| Column | Description |
|---|---|
| `sıra` | Optional sort/index |
| `ceza` | Penalty category name (`"KINAMA"`, `"OKULDAN KISA SÜRELİ UZAKLAŞTIRMA"`, `"OKUL DEĞİŞTİRME"`, `"ÖRGÜN EĞİTİM DIŞINA ÇIKARMA"`, `"UYARMA"`) |
| `madde` | Regulation clause code (e.g. `"164.1.a"`, `"55.a.1"`) |
| `açıklama` | Full text of the offense clause |
| `okultürü` | `"LİSE"` (74 rows) or `"ORTAOKUL"` (45 rows) |

*Note*: If `cezamaddeleri` is empty, `dsp_Modül.CezaMaddeleriKaydet()` seeds all 119 rows automatically.

#### 3. `Ögrenciler` (Student Master Roster per Academic Year)
| Column | Description |
|---|---|
| `sıra` | Sequential index |
| `ögretimyılı` | Academic year (`dsp_Modül.ÖgretimYılı`) |
| `okulno` | Student school number (unique key per `ögretimyılı`) |
| `adsoyad` | Student full name |
| `sınıf` | Class / section name (e.g. `"9. Sınıf / A Şubesi"` or `"9/A"`) |
| `cinsiyet` | Gender (`"KIZ"` or `"ERKEK"`) |
| `tcno` | 11-digit Turkish ID number |
| `dogumyeri` | Birthplace |
| `dogumtarihi` | Birthdate (`dd.MM.yyyy`) |
| `babaadı` | Father's name |
| `anneadı` | Mother's name |
| `veliadı` | Guardian's full name |
| `velitelno` | Guardian's phone number |
| `veliadres` | Guardian's address |
| `gündüzlü` | Boarding status (`"GÜNDÜZLÜ"` or `"PANSİYON / YATILI"`) |
| `davranışpuanı` | Current behavior score (defaults to `"100"` on creation/import) |
| `resimyolu` | Relative or absolute path to student photo (`resimler\{okulno}.jpg`) |

*(Note: `Ögrenciler` in `dsp.mdb` also has legacy columns mirroring `OlayÖgrenciler` from an older single-table schema, plus Ortaokul parent columns `anneyas`, `anneegitim`, `annehayatta`, `anneöz`, `annebirlikte`, `babayas`, `babaegitim`, `babahayatta`, `babaöz`, `bababirlikte`, `kardes`, which are read/written in `Frm04OlayÖgrenciEkleOrta`.)*

#### 4. `Olaylar` (Disciplinary Incidents)
| Column | Description |
|---|---|
| `sıra` | Sequential index |
| `ögretimyılı` | Academic year (`dsp_Modül.ÖgretimYılı`) |
| `olaykodu` | Auto-incremented numeric string per academic year (`"1"`, `"2"`, ...) |
| `olayadı` | Short title / subject of the incident |
| `yer` | Location where the incident occurred |
| `tarih` | Date of incident (`dd.MM.yyyy`) |
| `saat` | Time of incident |
| `kişi` | Person who reported / filed the petition (`Dilekçeyi Veren`) |
| `tarihsayı` | Date and reference number of the petition (`Dilekçe Tarih ve Sayısı`) |
| `özet` | Summary description of the incident (`Olayın Özeti`) |
| `dvgizle` | Visibility flag (unused/reserved) |

#### 5. `OlayÖgrenciler` (Students Involved in Incidents & Their Verdicts)
| Column | Description |
|---|---|
| `sıra` | Row ID |
| `ögretimyılı` | Academic year |
| `olaykodu` | Foreign key to `Olaylar.olaykodu` |
| `olayadı` | Copy of `Olaylar.olayadı` |
| `durumu` | Role in incident: `"SUÇLANAN ÖĞRENCİ"`, `"MAĞDUR / ŞİKAYETÇİ ÖĞRENCİ"`, or `"BİLGİSİNE BAŞVURULACAK ÖĞRENCİ"` |
| `tcno`, `adsoyad`, `okulno`, `sınıf`, `dogumtarihi`, `dogumyeri`, `babaadı`, `anneadı`, `cinsiyet`, `veliadı`, `velitelno`, `veliadres`, `gündüzlü` | Student identity & guardian fields copied/updated from `Ögrenciler` |
| `basaridurumu` | Academic standing (`"İYİ"`, `"ORTA"`, `"ZAYIF"`) |
| `saglıkdurumu` | Health status (`"NORMAL"`, `"RUHSAL HASTALIĞI VAR"`, `"BEDENSEL HASTALIĞI VAR"`, `"SÜREĞEN HASTALIĞI VAR"`) |
| `ekonomikdurum` | Family economic status (`"İYİ"`, `"NORMAL"`, `"ÇOK İYİ"`, `"ZAYIF"`, `"ÇOK ZAYIF"`) |
| `nerdekalıyor` | Living arrangement (`"AİLESİNİN YANINDA KALIYOR"`, etc.) |
| `annebabasag` | Parents alive status (`"ANNE BABA SAĞ"`, etc.) |
| `annebabaöz` | Parents biological/step status (`"ANNE VE BABA ÖZ"`, etc.) |
| `nerdeokuyor` | Commute status (`"AİLESİNİN YANINDA OKUYOR"`, etc.) |
| `büyüdügücevre` | Environment raised in (`"ŞEHİR MERKEZİ"`, `"İLÇE MERKEZİ"`, `"KASABA"`, `"KÖY"`) |
| `cezaaldımı` | Prior penalty status (`"CEZA ALMADI"` / `"CEZA ALDI"` or `"YAPTIRIM UYGULANMADI"` / `"YAPTIRIM UYGULANDI"`) |
| `okul` | Previous school attended (`Bitirdiği Okul` in Ortaokul) |
| `ifademetni` | Custom statement prompt text (`Frm05ifadeMetni`) |
| `yaziliuyari` | `"ALDI"` if given a administrative written warning (`Frm06YaziliUyari`) |
| `yaziliuyaritarihi` | Date of written warning |
| `yaziliuyaridavranis` | Offense behavior description for written warning |
| `yaziliuyarisayisi` | Official document number of written warning |
| `yaziliuyariceza` | Penalty category referenced in written warning |
| `yaziliuyarimadde` | Clause code referenced in written warning |
| `cezayeri` | Verdict Row 2: Place of offense (defaults to `Olaylar.yer`) |
| `cezatarihi` | Verdict Row 3: Date & time of offense (defaults to `"{Olaylar.tarih} - {Olaylar.saat}"`) |
| `cezacesidi` | Verdict Row 4: Offense category |
| `cezanedeni` | Verdict Row 5: Reason for penalty (defaults to `Olaylar.özet`) |
| `olayözet` | Verdict Row 6: Summary of student's defense (defaults to `"SUÇUNU KABUL EDİYOR"`) |
| `tanıkifadesi` | Verdict Row 7: Summary of witness statements |
| `deliller` | Verdict Row 8: Evidence (defaults to `"{Olaylar.kişi} nın {Olaylar.tarihsayı} tarih ve sayılı dilekçesi"`) |
| `hafifletennedenlar` | Verdict Row 9: Mitigating factors |
| `siddetlendirennedenler` | Verdict Row 10: Aggravating factors (auto-populated with prior penalties) |
| `kurulkanaati` | Verdict Row 11: Board's opinion (`"SUÇLU OLDUĞU YÖNÜNDE"` / `"SUÇSUZ OLDUĞU YÖNÜNDE"` or Ortaokul equivalents) |
| `yönmaddesi` | Verdict Row 12: Regulation clause code (e.g. `"164.1.a"`) |
| `yönceza` | Verdict Row 13: Regulation penalty corresponding to `yönmaddesi` |
| `yönaçıklama` | Verdict Row 14: Regulation clause description corresponding to `yönmaddesi` |
| `verilenceza` | Verdict Row 15: Actual penalty/sanction imposed by the board |
| `cezagünsayısı` | Verdict Row 16 (Lise only): Suspension duration (`"1 GÜN"` .. `"5 GÜN"`) |
| `cezauygulamatarihi` | Verdict Row 17 (Lise only): Suspension implementation dates |
| `inddavranışpuanı` | Verdict Row 18 (Lise) / Manual entry (`Frm09`): Behavior points deducted (`"10"`, `"20"`, `"40"`, `"80"`, or `"0"`) |
| `davranışpuanı` | Verdict Row 19 (Lise): Remaining behavior points after deduction |
| `karar` | Verdict Row 20 (Lise) / Row 16 (Ortaokul): Full text of the Board Verdict |
| `kararno` | Verdict Number (`"1"`, `"2"`, ... sequential per `ögretimyılı`) |
| `karartarihi` | Verdict Date (`dd.MM.yyyy`) |
| `iadedurumu` | End-of-year restoration status (`"İADE EDİLDİ"` / `"İADE EDİLMEDİ"` or `"YAPTIRIM KALDIRILDI"` / `"YAPTIRIM KALDIRILMADI"`) |
| `iadetarih` | Date of end-of-year point restoration / sanction removal |

#### 6. `gündemmaddeleri` (Meeting Agenda Items in `Frm06CagriPusulasi`)
| Column | Description |
|---|---|
| `sıra` | Agenda item number (`1`, `2`, ...) |
| `madde` | Agenda item text |

#### 7. `toplantitutanak` (Beginning/End of Year Meeting Decisions in `Frm19SeneBasiSonuKarar`)
| Column | Description |
|---|---|
| `sıra` | Row ID |
| `ögretimyılı` | Academic year |
| `toplantiadi` | `"SENE BAŞI"` or `"SENE SONU"` |
| `toplantitarihi` | Meeting date (`dd.MM.yyyy`) |
| `toplantisaati` | Meeting time (`HH:mm`) |
| `toplantiyeri` | Meeting location |
| `toplantibaskani` | Meeting chair name |
| `kararno` | Decision number (`1`, `2`, ...) |
| `karar` | Decision text (can include multi-line formatted tables of event/reward/penalty counts) |
| `açıklama` | Unused/empty string |

#### 8. `ödülcezasayilari` (Reward & Penalty Counts per Semester in `Frm19SeneBasiSonuKarar`)
| Column | Description |
|---|---|
| `sıra` | Row ID |
| `ögretimyılı` | Academic year |
| `toplantiadi` | `"SENE BAŞI"` or `"SENE SONU"` |
| `toplantitarihi`, `toplantisaati`, `toplantiyeri`, `toplantibaskani` | Meeting metadata |
| `ödülceza` | `"ödül"` or `"ceza"` |
| `ödülcezaadi` | Name of reward or penalty |
| `dönem1` | 1st Semester count |
| `dönem2` | 2nd Semester count |
| `toplam` | Total count (`dönem1 + dönem2`) |

### 2.4 Turkish Case Conversion Helpers (`dsp_Modül.cs`)
- `dsp_Modül.TurkceYap(deger)`: Converts strings to uppercase while preserving Turkish dotted/dotless I (`i -> İ`, `ı -> I`, `ğ -> Ğ`, `ü -> Ü`, `ş -> Ş`, `ö -> Ö`, `ç -> Ç`) before calling `.ToUpper()`.
- `dsp_Modül.IlkHarfBuyukYap(metin)`: Splits by space `' '`, capitalizes the first character of each word using Turkish rules (`i -> İ`, `ı -> I`), and lowercases the rest (`I -> ı`, `İ -> i`).

---

## 3. Screen-by-Screen (Form/Tab) Functional Specification

### 3.1 `AnaFrm.cs` — Main MDI Container & Ribbon/Toolbar
- **Startup sequence (`AnaFrm_Load`)**:
  1. Calls `dsp_Modül.Baglan()` and `dsp_Modül.TablolarınKontrolü()`.
  2. Checks if `cezamaddeleri` has rows; if empty, calls `dsp_Modül.CezaMaddeleriKaydet()`.
  3. Calls `dsp_Modül.KurumBilgiAl()` to load `Genelbilgiler` into global variables.
  4. If `dsp_Modül.OkulTürü == ""`, opens `Frm000OkulTürü` as a modal dialog so the user chooses `"LİSE"` or `"ORTAOKUL"`.
  5. Opens `Frm00Giris` (Dashboard) inside `Panel1`.
  6. Runs `Module_Pasword.DemoKontrol()`. If unlicensed, appends `"   - DEMO"` to the window title; if licensed, appends `"  -- (LİSANSLI) --"`.
- **School Type Adaptation**:
  - When `OkulTürü == "LİSE"`:
    - Window title: `"{ÖgretimYılı} - ÖDÜL VE DİSİPLİN KURULU PROGRAMI - LİSE   ({OkulAdı})"`
    - `BtnDavranisNotu` (`DAVRANIŞ PUANI İADE ÇİZELGESİ`) is visible; `BtnYaptırımKaldirma` is hidden.
  - When `OkulTürü == "ORTAOKUL"`:
    - Window title: `"{ÖgretimYılı} - ÖĞRENCİ DAVRANIŞLARINI DEĞERLENDİRME KURULU PROGRAMI - ORTAOKUL   ({OkulAdı})"`
    - `BtnYaptırımKaldirma` (`YAPTIRIM KARARI KALDIRMA`) is visible; `BtnDavranisNotu` is hidden.
    - Note: Clicking `BtnYaziliUyari` in `ORTAOKUL` mode shows an alert: `"Ortaokullarda yazılı uyarı uygulaması yoktur."` and does not open `Frm06YaziliUyari`.
- **Top Navigation Bar Buttons (`GroupBox1`)**:
  1. `BtnSabitBilgiler` ("GENEL BİLGİLER") -> `Frm01Genelbilgiler`
  2. `BtnÖgrenciBilgileri` ("ÖĞRENCİ BİLGİLERİ") -> `Frm02Ögrenciler`
  3. `BtnOlayKayit` ("OLAY GİRİŞİ") -> `Frm03OlayGirisi`
  4. `BtnOlayÖgrenci` ("OLAYA ÖĞRENCİ EKLEME") -> `Frm04OlayögrenciEkle` (if LİSE) or `Frm04OlayÖgrenciEkleOrta` (if ORTAOKUL)
  5. `BtnCagriPusulalari` ("ÇAĞRI PUSULALARI") -> `Frm06CagriPusulasi`
  6. `BtnifadeTutanak` ("İFADE VE BİLGİ TOPLAMA TUTANAKLARI") -> `Frm05ifadeTutanaklari`
  7. `BtnYaziliUyari` ("YAZILI UYARI") -> `Frm06YaziliUyari` (LİSE only)
  8. `BtnKarar` ("KURUL KARARI") -> `Frm07KararLise2` (if LİSE) or `Frm07KararOrta2` (if ORTAOKUL)
  9. `BtnCezateblig` ("CEZA TEBLİGATI") -> `Frm08CezaTabligati`
  10. `BtnDavranisNotu` ("DAVRANIŞ PUANI İADE ÇİZELGESİ") -> `Frm10DavranisPuaniCizelge` (LİSE only)
  11. `BtnYaptırımKaldirma` ("YAPTIRIM KARARI KALDIRMA") -> `Frm10yaptırımSilme` (ORTAOKUL only)
  12. `BtnYönetmelik` ("YÖNETMELİK MADDELERİ") -> `Frm13YönetmelikMaddeleri`
  13. `BtnKararlar` ("TÜM KARARLAR") -> `Frm07Kararlar`
  14. `BtnCezaAlanlar` ("CEZA ALAN ÖĞRENCİLER") -> `Frm12CezaAlanlar`
  15. `Btnistatistik` ("İSTATİSTİKLER") -> `Frm11Cezaistatistik`
  16. `BtnSeneBasiSonuKarar` ("SENE BAŞI/SONU TOPLANTI TUTANAĞI") -> `Frm19SeneBasiSonuKarar`
  17. `BtnDiziPusulasi` ("DİZİ PUSULASI") -> `Frm14Dizipusulasi`
  18. `BtnYedekleme` ("YEDEKLE") -> `Frm15Yedekleme` (Modal)
  19. `BtnGeriYükle` ("GERİ YÜKLE") -> `Frm16GeriYükleme` (Modal)
  20. `BtnHakkinda` ("LİSANS") -> `Frm17Hakkinda` (Modal)
  21. `BtnKapat` ("ÇIKIŞ") -> Closes `AnaFrm`.

### 3.2 `Frm000OkulTürü.cs` — School Type Selector
- Modal dialog with two radio buttons: `RBtnLise` (`"LİSE"`) and `RBtnOrtaokul` (`"ORTAOKUL"`).
- Saves selection to `Genelbilgiler.okultürü`.

### 3.3 `Frm00Giris.cs` — Home Dashboard
- Displays school info (`OkulAdı`, `ÖgretimYılı`) and quick-access tiles for the main workflow steps.

### 3.4 `Frm01Genelbilgiler.cs` — School & Board Settings
- Edits the single record in `Genelbilgiler`:
  - **Left Panel**: `CBoxOkulTürü` (`LİSE` / `ORTAOKUL`), `CBoxÖgretimYılı` (`2015 - 2016` .. `2030 - 2031`), `TxtKurumKodu`, `Txtİl` (auto-uppercased via `TurkceYap`), `Txtİlçe` (auto-uppercased via `TurkceYap`), `TxtOkulAdı` (title-cased via `IlkHarfBuyukYap`), `TxtMüdür` (title-cased), `TxtResmiYazıKodu`, `TxtAdres`, `TxtTelNo`, `TxtFaxNo`.
  - **Middle Panel**:
    - `TxtAntet`: Official multi-line document header (`T.C.\n{İLÇE} KAYMAKAMLIĞI\n{OkulAdı} Müdürlüğü`).
    - **Asil Üyeler (Primary Board Members)**: `TxtBaskan` + `TxtBaskanÜnvan` (default `"Müdür Yardımcısı"`), `TxtÜye1..4` + `TxtÜye1..4Ünvan` (default `"Öğretmen"` for 1–2, `"Onur Kurulu 2. Başkanı"` for 3 in Lise or `"Öğretmen"` in Ortaokul, `"Öğrenci Velisi"` for 4).
    - **Yedek Üyeler (Substitute Board Members)**: `TxtYedekBaskan` + `TxtYedekBaskanÜnvan`, `TxtYedekÜye1..4` + `TxtYedekÜye1..4Ünvan`.
  - **Right Panel**: Displays the statutory board composition rule (`MADDE 185` for Lise from `Label37.Text` or `MADDE 57` for Ortaokul from `Label11.Text`).

### 3.5 `Frm02Ögrenciler.cs` — Student Roster Management
- Displays all students in `Ögrenciler` for the active `ÖgretimYılı`.
- **Columns in `GridÖgrenciler`**:
  - `0`: `SIRA` (Read-only row index)
  - `1`: `NUMARASI` (`okulno`)
  - `2`: `ADI SOYADI` (`adsoyad`)
  - `3`: `SINIFI` (`sınıf`)
  - `4`: `CİNSİYETİ` (`cinsiyet` — ComboBox: `KIZ`, `ERKEK`)
  - `5`: `TC KİMLİK NO` (`tcno`)
  - `6`: `DOĞUM YERİ` (`dogumyeri`)
  - `7`: `DOĞUM TARİHİ` (`dogumtarihi`)
  - `8`: `BABA ADI` (`babaadı`)
  - `9`: `ANNE ADI` (`anneadı`)
  - `10`: `VELİ ADI SOYADI` (`veliadı`)
  - `11`: `VELİ TEL NO` (`velitelno`)
  - `12`: `VELİ ADRESİ` (`veliadres`)
  - `13`: `YATILILIK DURUMU` (`gündüzlü` — ComboBox: `GÜNDÜZLÜ`, `PANSİYON / YATILI`)
  - `14`: `DAVRANIŞ PUANI` (`davranışpuanı`, default `"100"`)
  - `15`: `RESİM YOLU` (`resimyolu`, hidden column)
- **Features**:
  - Filter by class (`CBoxSınıf`), search by Student No (`TxtNoArama`) or Name (`TxtAdArama`).
  - Double-clicking the `PictureBox` (`PBoxResim`) opens a file dialog to select a student photo, resizes it to `133x175`, saves it as `{StartupPath}\resimler\{okulno}.jpg`, and updates `Ögrenciler.resimyolu`.
  - Toolbar buttons launch:
    - `BtnExcel` -> `Frm02Excel` (Import students from e-Okul / e-Mesem Excel)
    - `BtnResimEkle` -> `Frm02ÖgrenciResimEkle` (Import student photos & data from e-Okul Excel / Word)
    - `BtnSinifAtlatma` -> `Frm02SinifAtlatma` (Copy/promote students from a previous academic year)
    - `BtnDavranisPuani` -> `Frm09DavranisPuaniGirisi` (Manual behavior point deduction entry)

### 3.6 `Frm02SinifAtlatma.cs` — Academic Year Class Promotion
- Allows selecting a source academic year (`CBoxEskiYıl`, e.g. `2024 - 2025`), source class (`CBoxEskiSınıf`), and target class (`TxtYeniSınıf`) in the current `dsp_Modül.ÖgretimYılı`.
- Copies all students from the source year/class into the active `ÖgretimYılı` with the new class name and resets `davranışpuanı` to `"100"`, skipping any student whose `okulno` already exists in the active year.

### 3.7 `Frm03OlayGirisi.cs` — Incident Entry
- Manages `Olaylar` for the active `ÖgretimYılı`.
- **Auto-Increment `OlayKoduAl()`**: Queries `SELECT * FROM Olaylar WHERE ögretimyılı='{ÖgretimYılı}'`, finds the maximum numeric `Val(olaykodu)`, and returns `max + 1` (starts at `1`).
- **Grid Columns**:
  - `0`: `OLAY KODU` (`olaykodu`, ReadOnly)
  - `1`: `OLAY ADI / KONUSU` (`olayadı`)
  - `2`: `OLAY YERİ` (`yer`)
  - `3`: `OLAY TARİHİ` (`tarih`, picked via floating `DateTimePicker`)
  - `4`: `OLAY SAATİ` (`saat`)
  - `5`: `DİLEKÇEYİ VEREN / BİLDİREN` (`kişi`)
  - `6`: `DİLEKÇE TARİH VE SAYISI` (`tarihsayı`)
  - `7`: `OLAYIN ÖZETİ` (`özet`)
- Deleting an incident deletes its row from `Olaylar` AND cascades (`DELETE FROM OlayÖgrenciler WHERE olaykodu='{olaykodu}' AND ögretimyılı='{ÖgretimYılı}'`).

### 3.8 `Frm04OlayögrenciEkle.cs` (LİSE) & `Frm04OlayÖgrenciEkleOrta.cs` (ORTAOKUL) — Adding Students to Incidents
- Left side lists incidents from `Olaylar` for `ÖgretimYılı`. Selecting an incident loads its linked students from `OlayÖgrenciler`.
- Clicking `BtnÖgrenciEkle` opens `Frm04OgrenciAl`:
  - Displays all students from `Ögrenciler` (searchable by name/number/class) with a role selector (`CBoxDurumu`: `"SUÇLANAN ÖĞRENCİ"`, `"MAĞDUR / ŞİKAYETÇİ ÖĞRENCİ"`, `"BİLGİSİNE BAŞVURULACAK ÖĞRENCİ"`).
  - Alternatively allows typing an external/unregistered student manually.
- **Student Socio-Economic & Family Profile Fields**:
  - **In `Frm04OlayögrenciEkle.cs` (LİSE)**:
    - `GridÖgrenciBilgileri` is a 2-column vertical property grid with 22 rows:
      - `0`: `Olaydaki Durumu` (ComboBox: `SUÇLANAN ÖĞRENCİ`, `MAĞDUR / ŞİKAYETÇİ ÖĞRENCİ`, `BİLGİSİNE BAŞVURULACAK ÖĞRENCİ`)
      - `1`: `Adı Soyadı`
      - `2`: `Sınıfı`
      - `3`: `Okul Numarası`
      - `4`: `Cinsiyeti` (`KIZ`, `ERKEK`)
      - `5`: `TC Kimlik Numarası`
      - `6`: `Doğum Yeri`
      - `7`: `Doğum Tarihi`
      - `8`: `Baba Adı`
      - `9`: `Anne Adı`
      - `10`: `Veli Adı Soyadı`
      - `11`: `Veli Telefon Numarası`
      - `12`: `Veli Adresi`
      - `13`: `Yatılılık Durumu` (`GÜNDÜZLÜ`, `PANSİYON / YATILI`)
      - `14`: `Başarı Durumu` (`İYİ`, `ORTA`, `ZAYIF`)
      - `15`: `Sağlık Durumu` (`NORMAL`, `RUHSAL HASTALIĞI VAR`, `BEDENSEL HASTALIĞI VAR`, `SÜREĞEN HASTALIĞI VAR`)
      - `16`: `Ailenin Ekonomik Durumu` (`İYİ`, `NORMAL`, `ÇOK İYİ`, `ZAYIF`, `ÇOK ZAYIF`)
      - `17`: `Öğrencinin Nerede Kaldığı` (`AİLESİNİN YANINDA KALIYOR`, `AKRABASININ YANINDA KALIYOR`, `DEVLET YURDUNDA KALIYOR`, `ÖZEL YURTTA KALIYOR`, `ARKADAŞLARIYLA EVDE KALIYOR`, `KENDİSİ AYRI EVDE KALIYOR`)
      - `18`: `Anne - Baba Hayatta mı?` (`ANNE BABA SAĞ`, `ANNE SAĞ BABA ÖLÜ`, `BABA SAĞ ANNE ÖLÜ`, `ANNE VE BABA ÖLÜ`)
      - `19`: `Anne - Baba Öz mü?` (`ANNE VE BABA ÖZ`, `ANNE ÖZ BABA ÜVEY`, `BABA ÖZ ANNE ÜVEY`, `ANNE VE BABA ÜVEY`)
      - `20`: `Ailesinin Yanında mı Okuyor?` (`AİLESİNİN YANINDA OKUYOR`, `GÜNÜBİRLİK GELİP GİDİYOR`, `AKRABASININ YANINDA OKUYOR`)
      - `21`: `Büyüdüğü Çevre` (`ŞEHİR MERKEZİ`, `İLÇE MERKEZİ`, `KASABA`, `KÖY`)
      - `22`: `Daha Önce Ceza Aldı mı?` (`CEZA ALMADI`, `CEZA ALDI`)
  - **In `Frm04OlayÖgrenciEkleOrta.cs` (ORTAOKUL)**:
    - Adds detailed parent comparison fields (stored in `Ögrenciler` and printed in `EK-10` Karar table):
      - `Bitirdiği Okul` (`okul`), `Kardeş Sayısı` (`kardes`)
      - **Anne**: `Anne Adı` (`anneadı`), `Yaşı` (`anneyas`), `Eğitim Durumu` (`anneegitim`), `Hayatta mı?` (`annehayatta`: `EVET`/`HAYIR`), `Öz mü?` (`anneöz`: `ÖZ`/`ÜVEY`), `Birlikte mi?` (`annebirlikte`: `EVET`/`HAYIR`)
      - **Baba**: `Baba Adı` (`babaadı`), `Yaşı` (`babayas`), `Eğitim Durumu` (`babaegitim`), `Hayatta mı?` (`babahayatta`: `EVET`/`HAYIR`), `Öz mü?` (`babaöz`: `ÖZ`/`ÜVEY`), `Birlikte mi?` (`bababirlikte`: `EVET`/`HAYIR`)
- Saving updates both `OlayÖgrenciler` (for that incident) and synchronizes shared demographic fields back to `Ögrenciler`.

### 3.9 `Frm06CagriPusulasi.cs` — Summons Slips & Board Meeting Notice
- Contains two tabs (`TabControl1`):
  1. **TabPage1 (`ÖĞRENCİ ÇAĞRI PUSULASI`)**:
     - Lists students involved in incidents (`GridÖgrenciler`) with a radio filter:
       - `RbtnGöster1`: `"TÜM ÖĞRENCİLER"`
       - `RbtnGöster2`: `"SUÇLANAN ÖĞRENCİ"`
       - `RbtnGöster3`: `"MAĞDUR / ŞİKAYETÇİ ÖĞRENCİ"`
       - `RbtnGöster4`: `"BİLGİSİNE BAŞVURULACAK ÖĞRENCİ"`
     - Inputs: `DTtarih` (Date), `TxtSaat` (Time), `TxtYer` (Place, e.g. `"Disiplin Kurulu Odası"`).
     - Buttons:
       - `BtnSavunmaCagri` ("SAVUNMA ÇAĞRI PUSULASI YAZDIR") -> Opens `Frm07Yazdir` (board member selector) then calls `YazÖgrenciSavunmaCagri()` (`ögrencisavunmacagripusulasi{N}.pdf`).
       - `BtnİfadeCagri` ("İFADE ÇAĞRI PUSULASI YAZDIR") -> Opens `Frm07Yazdir` then calls `YazÖgrenciİfadeCagri()` (`ögrenciifadecagripusulasi{N}.pdf`).
       - `BtnYazdir` ("ÖĞRENCİ LİSTESİ YAZDIR") -> Calls `YazdirListe()` (`cagri-pusulasi-liste{N}.pdf`).
  2. **TabPage2 (`KURUL ÜYELERİ TOPLANTI ÇAĞRI PUSULASI`)**:
     - Inputs: `TxtSayı` (defaults to `dsp_Modül.ResmiYazıKodu`), `TxtKonu` (defaults to `"{KurulAdı} Kurulu Toplantısı"`), `TxtTarih` (Document Date), plus meeting `DTtarih`, `TxtSaat`, `TxtYer` from Tab 1.
     - `GridGündem`: Editable list of meeting agenda items persisted to table `gündemmaddeleri` via `GündemKaydet()`.
     - Button `BtnÜyeCagri` ("ÜYE ÇAĞRI PUSULASI YAZDIR") -> Opens `Frm07Yazdir` then calls `YazÜyeCagri()` (`üyecagripusulasi{N}.pdf`).

### 3.10 `Frm05ifadeTutanaklari.cs` & `Frm05ifadeMetni.cs` — Statement, Defense, & Information Gathering Forms
- Lists incidents (`GridOlaylar`) and the students linked to the selected incident (`GridOlayÖgrenciler`).
- Clicking `"İFADE METNİ GİRİŞİ"` opens `Frm05ifadeMetni` to edit `OlayÖgrenciler.ifademetni` (defaults to `"{Olaylar.tarih} tarihinde {Olaylar.yer} nde meydana gelen {Olaylar.olayadı} olayı ile ilgili bildiklerinizi yazınız."`).
- Provides 14 PDF generation buttons (detailed in Section 6.1):
  1. `İfade Tutanağı` (`ifade-tutanagi{N}.pdf`)
  2. `Savunma Tutanağı` (`savunma-tutanagi{N}.pdf`)
  3. `Bilgi Toplama Formu - Olay` (`bilgitoplamaolay-tutanagi{N}.pdf`)
  4. `Bilgi Toplama Formu - Öğretmen` (`bilgitoplamaögretmen-tutanagi{N}.pdf`)
  5. `Bilgi Toplama Formu - Sınıf Rehber Öğretmeni` (`bilgitoplamasınıfögretmen-tutanagi{N}.pdf`)
  6. `Bilgi Toplama Formu - Veli` (`bilgitoplamaveli-tutanagi{N}.pdf`)
  7. `Bilgi Toplama Formu - Öğrenci` (`bilgitoplamaörenci-tutanagi{N}.pdf`)
  8. `Ek Süre Talep Formu` (`eksüretalep-form{N}.pdf` — uses inputs `TxtSüre0` for extra days and `TxtSüre1` for reason)
  9. `Şikayet Formu` (`sikayet-form{N}.pdf`)
  10. `Sözlü Uyarı Tutanağı` (`sözlü-uyari-tutanagi{N}.pdf`)
  11. `Öğrenci Sözleşmesi` (`ögrenci-sözlesmesi{N}.pdf`)
  12. `Veli Görüşme Tutanağı` (`veli-görüsme-tutanagi{N}.pdf`)
  13. `Veli Görüşme Formu` (`veli-görüsme-form{N}.pdf`)
  14. `Rehber Öğretmen Görüş Formu` (`rehber-ögretmen-görüs-form{N}.pdf`)

### 3.11 `Frm06YaziliUyari.cs` — Administrative Written Warning (LİSE Only)
- Used in LİSE under *MEB Ortaöğretim Kurumları Yönetmeliği Madde 163/2* when school management issues a written warning instead of referring the student to a disciplinary penalty.
- Lists incidents and students. Selecting a student loads their written warning fields from `OlayÖgrenciler`:
  - `TxtSayi`: Defaults to `"{ResmiYazıKodu}/{Genelbilgiler.uyarisayi + 1}"`
  - `DTTarih`: Warning date (`yaziliuyaritarihi`)
  - `TxtDavranis`: Behavior description (`yaziliuyaridavranis`, defaults to `Olaylar.özet`)
  - `CBoxCeza`: Penalty category (`yaziliuyariceza`)
  - `CBoxMadde`: Clause code (`yaziliuyarimadde`) -> auto-fills `TxtAciklama` from `cezamaddeleri`.
- Saving sets `OlayÖgrenciler.yaziliuyari = "ALDI"`, saves the warning fields, and increments `Genelbilgiler.uyarisayi`.
- Printing generates `yazili-uyari{N}.pdf`.

### 3.12 `Frm07KararLise2.cs` (LİSE) & `Frm07KararOrta2.cs` (ORTAOKUL) — Board Verdict Entry
- Lists incidents (`GridOlaylar`) and accused students (`GridOlayÖgrenciler` where `durumu = 'SUÇLANAN ÖĞRENCİ'`).
- Edits the 21-row (`Frm07KararLise2`) or 17-row (`Frm07KararOrta2`) verdict property grid (`GridÖgrenciBilgileri`), calculates behavior point deductions, auto-generates the formal verdict text (`KurulKararı()`), saves to `OlayÖgrenciler` and `Ögrenciler`, and prints `karar-tutanak{N}.pdf` (`EK-1` for LİSE, `EK-10` for ORTAOKUL).
- See **Section 5** for the complete calculation and text generation rules.

### 3.13 `Frm07Kararlar.cs` — All Verdicts Register
- Displays a read-only master list of all students who received a verdict (`kararno <> ''`) in the active `ÖgretimYılı`, sorted by `Val(kararno)`.

### 3.14 `Frm08CezaTabligati.cs` — Penalty Notification Slips (Student & Parent)
- Lists all students in `OlayÖgrenciler` for `ÖgretimYılı` who have a non-empty `verilenceza` and `kararno`.
- Two print actions:
  1. `BtnÖgrenciTebligat` ("ÖĞRENCİ KARAR TEBLİGAT PUSULASI") -> `ögrencitabligat{N}.pdf` (`YazÖgrenciTebligatPDF()`)
  2. `BtnVeliTebligat` ("VELİ CEZA BİLDİRİM YAZISI") -> `Velitabligat{N}.pdf` (`YazVeliTebligat()`)

### 3.15 `Frm09DavranisPuaniGirisi.cs` — Manual Behavior Score Entry
- Allows viewing and directly editing `inddavranışpuanı` and `davranışpuanı` for students in `OlayÖgrenciler` / `Ögrenciler`.

### 3.16 `Frm10DavranisPuaniCizelge.cs` (LİSE) & `Frm10yaptırımSilme.cs` (ORTAOKUL) — End-of-Year Score Restoration / Sanction Removal
- **LİSE (`Frm10DavranisPuaniCizelge.cs` — Madde 171/2)**:
  - Lists all students in `OlayÖgrenciler` for `ÖgretimYılı` who received a penalty (`verilenceza` in `KINAMA`, `OKULDAN KISA SÜRELİ UZAKLAŞTIRMA`, `OKUL DEĞİŞTİRME`, `ÖRGÜN EĞİTİM DIŞINA ÇIKARMA`).
  - Columns include a checkbox (`SEÇ`), student info, incident info, penalty, deducted points (`inddavranışpuanı`), current score (`davranışpuanı`), restoration date (`iadetarih`, picked via `DTiadeTarihi`), and restoration status (`iadedurumu`: ComboBox `"İADE EDİLDİ"` / `"İADE EDİLMEDİ"`).
  - Saving updates `OlayÖgrenciler.iadetarih` and `OlayÖgrenciler.iadedurumu`.
  - Printing generates landscape A4 PDF `davranispuaniiadecizelgesi{N}.pdf`.
- **ORTAOKUL (`Frm10yaptırımSilme.cs` — Madde 62/3)**:
  - Lists all students in `OlayÖgrenciler` for `ÖgretimYılı` who received a sanction (`UYARMA`, `KINAMA`, `OKUL DEĞİŞTİRME`).
  - Allows setting `iadetarih` and `iadedurumu` (`"YAPTIRIM KALDIRILDI"` / `"YAPTIRIM KALDIRILMADI"`).
  - Printing generates landscape A4 PDF `yaptirim-karari-kaldirma-raporu{N}.pdf`.

### 3.17 `Frm11Cezaistatistik.cs` — Disciplinary Statistics & Charts
- Three tabs with `DataGridView` + MSChart (`System.Windows.Forms.DataVisualization.Charting.Chart`) visualization (switchable between Column chart and Pie chart, 2D/3D):
  1. **Tab 1 (`VERİLEN CEZALARA GÖRE İSTATİSTİKLER`)**: Counts occurrences of each penalty in `OlayÖgrenciler` for `ÖgretimYılı` (note: for `"YAZILI OLARAK UYARMA"`, queries `yaziliuyari='ALDI'`). Prints `cezaistatistik{N}.pdf`.
  2. **Tab 2 (`OLAYLARA GÖRE İSTATİSTİKLER`)**: Lists each incident in `Olaylar` and counts how many students in `OlayÖgrenciler` were involved in that `olaykodu`. Prints `olayistatistik{N}.pdf`.
  3. **Tab 3 (`YÖNETMELİK MADDELERİNE GÖRE İSTATİSTİKLER`)**: Counts occurrences of each `yönmaddesi` in `OlayÖgrenciler` for the school type's clauses. Prints `yönetmelikistatistik{N}.pdf`.
- Users can choose Portrait (`DİKEY`) or Landscape (`YATAY`) and whether to embed the chart image (`CBoxGrafik`: `"GRAFİK EKLE"`) into the PDF table.

### 3.18 `Frm12CezaAlanlar.cs` — Penalized Students Report
- Two tabs filtering `OlayÖgrenciler` for `ÖgretimYılı`:
  1. **Tab 1 (`SINIFLARA GÖRE CEZA ALAN ÖĞRENCİLER`)**: Filter by class (`CBoxSınıf`, including `"TÜM SINIFLAR"`).
  2. **Tab 2 (`CEZALARA GÖRE ÖĞRENCİLER`)**: Filter by penalty type (`CBoxCeza`).
- Prints `sınıfcezaalanögrencilistesi{N}.pdf`.

### 3.19 `Frm13YönetmelikMaddeleri.cs` — Regulation Clauses Editor
- Two tabs (`LİSE` and `ORTAOKUL`) showing `cezamaddeleri` in editable grids (`CEZA`/`YAPTIRIM`, `MADDE NO`, `MADDE AÇIKLAMASI`).
- Supports adding rows (`SATIR EKLE`), deleting rows (`SATIR SİL`), searching by text (`TxtArama`), and saving (`KAYDET / DEĞİŞTİR`, which drops and recreates `cezamaddeleri` from both grids).

### 3.20 `Frm14Dizipusulasi.cs` — Investigation File Document Index (`Dizi Pusulası`)
- Dual-grid builder:
  - Left grid (`Grid1`): Pre-populated with 27 standard disciplinary file document types (see `yonetmelik.json -> document_templates.dizi_pusulasi.default_items`).
  - Right grid (`Grid2`): Documents selected for the current file (`SIRA`, `KONUSU`, `ADET / SAYI` defaulting to `"1"`).
  - Buttons move items between `Grid1` and `Grid2` (`EKLE >`, `< ÇIKAR`) and reorder rows in `Grid2` (`ÜSTE ÇIKAR`, `ALTA İNDİR`).
  - `YAZDIR` generates `{N}.pdf` (`DİZİ PUSULASI`).

### 3.21 `Frm15Yedekleme.cs` & `Frm16GeriYükleme.cs` — Backup & Restore
- **Backup (`Frm15Yedekleme`)**: Closes the ADODB connection (`dsp_Modül.BaglantıKapat()`), copies `dsp.mdb` to the user-selected `.mdb` path, and reconnects.
- **Restore (`Frm16GeriYükleme`)**: Closes the ADODB connection, deletes `dsp.mdb`, copies the user-selected backup `.mdb` to `dsp.mdb`, reconnects, and runs `dsp_Modül.TablolarınKontrolü()`.

### 3.22 `Frm19SeneBasiSonuKarar.cs` — Beginning & End of Year Board Meeting Minutes
- Manages meeting minutes for `"SENE BAŞI"` and `"SENE SONU"` meetings in `toplantitutanak` and `ödülcezasayilari`.
- **Left Panel**: `CBOXtopAdi` (`"SENE BAŞI"`, `"SENE SONU"`), `DTPtoptarihi`, `MTBOXtopsaati` (`00:00`), `TBOXtopyeri`, `CBoxTopBaskani` (contains `KurulBaşkanı` and `YedekKurulBaşkanı`).
- **Center Panel (`dgvKararlar`)**: Multi-line auto-sizing grid of decisions (`NO`, `KARAR`) with Undo stack (`BtnGeriAl`), insert row above/below, delete row, and move row up/down. Pressing `Shift+Enter` (`keyData == 131085`) inside a cell inserts a newline.
- **Right Panel (`İSTATİSTİKLER`)**:
  - **Ödül Sayıları (`GridÖdül`)**: 4 rows (`TEŞEKKÜR BELGESİ`, `TAKDİR BELGESİ`, `ÜSTÜN BAŞARI BELGESİ`, `ONUR BELGESİ`) with editable `1. DÖNEM` and `2. DÖNEM` counts and auto-summed `TOPLAM`. Clicking `"ÖDÜL SAYISINI RAPORA EKLE"` formats a fixed-width text table (`Sabitle(40)` + `3 x Sabitle(12)` with `'—'` separator lines) and appends it as a decision row in `dgvKararlar`.
  - **Olay Sayısı (`TxtOlaySayısı`)**: Counts `SELECT COUNT(*) FROM olaylar WHERE ögretimyılı='{ÖgretimYılı}'`. Clicking `"OLAY SAYISINI RAPORA EKLE"` appends `"Okul Müdürlüğünce, {KurulAdı} Kurulumuza toplam {TxtOlaySayısı} olay sevk edilmiştir."` to `dgvKararlar`.
  - **Ceza Sayıları (`GridCezalar`)**:
    - Automatically counts penalties by semester (`CezaÖgrenciSay`):
    - **Semester Cutoff Date Calculation**: Splits `dsp_Modül.ÖgretimYılı` by `'-'`, parses the second year `Y2`, and sets the 1st-semester cutoff date to **`January 30, Y2` (`new DateTime(year, 1, 30)`)**.
    - Verdicts with `karartarihi <= Jan 30, Y2` are counted in `1. DÖNEM`; verdicts with `karartarihi > Jan 30, Y2` are counted in `2. DÖNEM`. (For `"YAZILI OLARAK UYARMA"`, uses `yaziliuyari='ALDI'` and `yaziliuyaritarihi`).
    - Clicking `"CEZA SAYISINI RAPORA EKLE"` formats a fixed-width text table (`Sabitle(40)` + `3 x Sabitle(12)`) and appends it to `dgvKararlar`.
- Clicking `YAZDIR - PDF` generates `sene-basi-sonu-tutanak{N}.pdf` using both `Tahoma` and `Courier New` (`cour.ttf`) so the fixed-width tables align in the PDF.

---

## 4. e-Okul / e-Mesem Excel & Word Parsing Logic

The application has two distinct import mechanisms for student data:

### 4.1 Student Roster Import from e-Okul / e-Mesem Excel (`Frm02Excel.cs`)

#### Supported Formats & OLEDB Connection
- Accepts `.xls` and `.xlsx` files via `OpenFileDialog`.
- Opens the workbook using `System.Data.OleDb.OleDbConnection`:
  ```text
  Provider=Microsoft.ACE.OLEDB.12.0;Data Source={dosya};Extended Properties="Excel 12.0 Xml;HDR=NO";
  ```
- **Worksheet Discovery (`SayfaAdiBul`)**:
  - Checks `GetOleDbSchemaTable(OleDbSchemaGuid.Tables, null)` for any sheet in `["Sayfa1$", "Sheet1$", "RAPOR$", "Rapor1$"]` (case-insensitive).
  - If none of those 4 names match, falls back to the first `TABLE_NAME` in the schema.
  - Reads the entire sheet via `SELECT * FROM [{sheetName}]` into a `DataTable` (`HDR=NO`, so row 0 of the sheet is `DataTable.Rows[0]`).

#### Supported e-Okul Report Types
As instructed in the UI (`Frm02Excel`):
1. **e-Okul**: *Öğrenci İşlemleri -> Raporlar -> Öğrenci Künye Defteri -> Excel (Sadece Veri)* (or Class List exports).
2. **e-Mesem**: Student list exports where Col 1 is Student No, Col 2 is Full Name (`Ad Soyad` combined), and Col 7 is Class/Section.

#### Row-by-Row Parsing Algorithm (`BtnExcelYükle_Click`)
The parser iterates `i = 0` to `DataTable.Rows.Count - 1`:

1. **Class Header Banner Detection (Standard e-Okul Class List)**:
   - Inspects ` cell0 = Trim(Rows[i][0]) `.
   - If `cell0` contains `"Sınıfı / Şubesi"` (or starts with a class name banner where `Rows[i][0]` has `"Sınıf"`):
     - Extracts the current class name `SınıfAdı` (e.g. `"9. Sınıf / A Şubesi"`).
2. **Student Row Detection & Column Mapping**:
   - Checks if `Val(Rows[i][1]) > 0` (Column index `1` is numeric Student Number `okulno`) OR `Val(Rows[i][0]) > 0`:
   - **Layout A — Standard e-Okul "Öğrenci Sınıf Listesi" (Class Banner above rows)**:
     - Col `1`: `OkulNo` (numeric)
     - Col `2` or `3`: First Name (`Adı`)
     - Col `4` or `6`: Last Name (`Soyadı`)
     - Col `8` or `9`: Gender (`Cinsiyeti`: `"K"`/`"KIZ"` -> `"KIZ"`, `"E"`/`"ERKEK"` -> `"ERKEK"`)
     - Class is taken from the most recently seen Class Header Banner (`SınıfAdı`).
   - **Layout B — e-Okul "Öğrenci Künye Defteri" (Excel Sadece Veri)**:
     - Contains full demographic columns on each row:
       - `OkulNo`, `AdSoyad`, `Sınıf`, `TCKimlikNo`, `DoğumYeri`, `DoğumTarihi`, `BabaAdı`, `AnneAdı`, `Cinsiyet`.
   - **Layout C — e-Mesem Format**:
     - Triggered when Col `1` is numeric (`OkulNo`), Col `2` contains the full name (`Ad Soyad`), Col `4` is empty, and Col `7` contains the class/branch name (`Sınıf`).

#### Three Import Modes in `Frm02Excel`
Before saving to `Ögrenciler`, the user chooses one of three radio buttons:
1. **`RBtnYeniKayit` ("SADECE YENİ ÖĞRENCİLERİ EKLE")**:
   - Checks if `okulno` already exists in `Ögrenciler` for `dsp_Modül.ÖgretimYılı`.
   - Only inserts students whose `okulno` is not yet in the database; sets `davranışpuanı = "100"` and `gündüzlü = "GÜNDÜZLÜ"`.
2. **`RBtnTümünüSil` ("TÜM ÖĞRENCİLERİ SİL VE YENİDEN YÜKLE")**:
   - Executes `DELETE FROM Ögrenciler WHERE ögretimyılı='{dsp_Modül.ÖgretimYılı}'` and inserts all parsed rows with `davranışpuanı = "100"`.
3. **`RBtnSınıfGüncelle` ("ÖĞRENCİ SINIFLARINI GÜNCELLE")**:
   - Matches existing students by `okulno` and `ögretimyılı` and updates their `sınıf` field (and inserts any new students).

---

### 4.2 Resimli Öğrenci Listesi (Photo + Roster) Import (`Frm02ÖgrenciResimEkle.cs`)

This form imports student photos (`resimler\{OkulNo}.jpg`) AND optionally inserts/updates student records (`OkulNo`, `AdSoyad`, `Sınıf`) from **e-Okul -> Öğrenci İşlemleri -> Raporlar -> Resimli Öğrenci Listesi -> Excel (sadece veri)** (`.xls` or `.xlsx`) or Word (`.docx`).

#### Mode 1: `.xls` (Legacy BIFF8 Excel) via Office Interop + Clipboard (`ResimAlXLS`)
1. Launches `Microsoft.Office.Interop.Excel.Application` invisibly and opens the `.xls` workbook.
2. Scans cells in `Range["A1:Z500"]`:
   - Any cell containing `"Sınıfı"` and `"Şubesi"` sets the current `Sınıf` string.
   - Any cell whose text is `"Numarası:"` or `"Okul No:"` has its adjacent cell read as `OkulNo`, and the cell below/next to `"Adı Soyadı:"` read as `AdSoyad`.
3. Iterates through `Worksheet.Shapes`:
   - Calls `shape.CopyPicture()`, retrieves the `System.Drawing.Image` from `Clipboard.GetImage()`, resizes it to **`266 x 350`** pixels (`new Bitmap(img, 266, 350)`), and saves it to `{StartupPath}\resimler\{OkulNo}.jpg` in order of appearance.

#### Mode 2: `.xlsx` (OpenXML Excel) via Zip/XML Drawing Anchor Inspection (`ResimAlXLSX`)
1. Opens the `.xlsx` file as a `System.IO.Compression.ZipArchive` AND opens its cells via OLEDB/EPPlus (`OfficeOpenXml.ExcelPackage`).
2. Reads `xl/drawings/drawing1.xml` and `xl/drawings/_rels/drawing1.xml.rels`:
   - For each `<xdr:twoCellAnchor>` or `<xdr:oneCellAnchor>`, extracts:
     - Top-left anchor `(fromRow, fromCol)` (`<xdr:from><xdr:row>`, `<xdr:col>`)
     - Relationship ID `r:embed="rId..."` on `<a:blip>` -> resolves target media path `xl/media/image{N}.{ext}` from `drawing1.xml.rels`.
3. **Coordinate Matching for Student Metadata**:
   - For a photo anchored at `(fromRow, fromCol)`:
     - Searches upward from `fromRow` for the nearest row containing `"Sınıfı / Şubesi"` (or `"Sınıfı"`) to determine `Sınıf`.
     - Searches the cells immediately to the right or below `(fromRow, fromCol)` (within rows `fromRow .. fromRow + 6`, cols `fromCol .. fromCol + 5`) for the student's `OkulNo` (numeric string) and `AdSoyad`.
4. Extracts the image stream from `xl/media/...`, resizes it to **`266 x 350`** JPEG, saves it to `{StartupPath}\resimler\{OkulNo}.jpg`, and updates `Ögrenciler.resimyolu` (and inserts the student into `Ögrenciler` if not already present).

#### Mode 3: `.docx` (Word Document) via Zip Media Extraction (`ResimAlDOCX`)
- Opens the `.docx` archive, extracts images from `word/media/`, matches them against student numbers parsed from the document tables, and saves `resimler\{OkulNo}.jpg`.

---

## 5. Verdict, Penalty, & Behavior Point Calculation Logic

### 5.1 Accused Student Filtering
- Only students in `OlayÖgrenciler` with `durumu == "SUÇLANAN ÖĞRENCİ"` are loaded into `Frm07KararLise2` and `Frm07KararOrta2`.

### 5.2 Board Member Conflict-of-Interest Substitution
In both `Frm07KararLise2` and `Frm07KararOrta2` (`OlayBilgisiAl` / `KurulÜyeleriKontrol`):
1. First, primary board members (`KurulBaşkanı`, `Üye1`..`Üye4`) and their titles are loaded from `Genelbilgiler`.
2. Next, the complainant's name (`DilekçeyiVeren = Olaylar.kişi`) is compared (case-insensitive via `TurkceYap`) against each primary board member:
   - If `DilekçeyiVeren == KurulBaşkanı` -> replaced by `YedekKurulBaşkanı` (and `YedekKurulBaşkanıünvan`).
   - If `DilekçeyiVeren == Üye1` -> replaced by `YedekÜye1` (and `YedekÜye1ünvan`).
   - If `DilekçeyiVeren == Üye2` -> replaced by `YedekÜye2` (and `YedekÜye2ünvan`).
   - If `DilekçeyiVeren == Üye3` -> replaced by `YedekÜye3` (and `YedekÜye3ünvan`).
   - If `DilekçeyiVeren == Üye4` -> replaced by `YedekÜye4` (and `YedekÜye4ünvan`).
3. Furthermore, when the user clicks `YAZDIR` (`BtnYazdir`), `Frm07Yazdir` pops up before PDF generation, displaying a dropdown for each of the 5 board positions (`Başkan`, `Üye 1`..`4`) pre-populated with `[Asil Member, Yedek Member]` so the user can manually override who attended the meeting.

### 5.3 Sequential Verdict Numbering (`KararNoAl()`)
- Queries `SELECT kararno FROM OlayÖgrenciler WHERE ögretimyılı='{dsp_Modül.ÖgretimYılı}'`.
- Computes `EnBüyük = MAX(Conversion.Val(kararno))`.
- Next verdict number is `EnBüyük + 1` (starts at `1` at the beginning of each academic year).

### 5.4 Automatic Prior Penalties Lookup (`AldıgiCezalar(adsoyad)`)
- When loading an accused student who does not yet have a saved verdict (`kararno == ""`):
- Queries `SELECT * FROM OlayÖgrenciler WHERE adsoyad='{adsoyad}' AND kararno<>''`.
- Concatenates every prior record into a single string:
  ```text
  "{ögretimyılı} Öğretim Yılı {verilenceza} Cezası "
  ```
- Places this string into **Row 10 (`Ceza Şiddetlendiren Nedenler` / `siddetlendirennedenler`)**.

### 5.5 LİSE Verdict Grid (`Frm07KararLise2.cs` — 21 Rows, 0-indexed `0..20`)

| Row Index | Label (`Column1`) | DB Column (`OlayÖgrenciler`) | Control / Behavior |
|---|---|---|---|
| `0` | `Karar Numarası` | `kararno` | Auto-filled by `KararNoAl()` if empty |
| `1` | `Karar Tarihi` | `karartarihi` | `DateTimePicker` (defaults to `Now.ToShortDateString()`) |
| `2` | `Suç Yeri` | `cezayeri` | Defaults to `Olaylar.yer` |
| `3` | `Suç Tarihi ve Saati` | `cezatarihi` | Defaults to `"{Olaylar.tarih} - {Olaylar.saat}"` |
| `4` | `Suç Çeşidi` | `cezacesidi` | ComboBox: `KINAMA CEZASINI GEREKTİREN DAVRANIŞ`, `OKULDAN KISA SÜRELİ UZAKLAŞTIRMA CEZASINI GEREKTİREN DAVRANIŞ`, `OKUL DEĞİŞTİRME CEZASINI GEREKTİREN DAVRANIŞ`, `ÖRGÜN EĞİTİM DIŞINA ÇIKARMA CEZASINI GEREKTİREN DAVRANIŞ` |
| `5` | `Suçun Nedeni` | `cezanedeni` | Defaults to `Olaylar.özet` |
| `6` | `Olayın Özeti` | `olayözet` | Defaults to `"SUÇUNU KABUL EDİYOR"` |
| `7` | `Tanıkların İfadesinin Özeti` | `tanıkifadesi` | Free text |
| `8` | `Suçu Kanıtlayan Deliller` | `deliller` | Defaults to `"{Olaylar.kişi} nın {Olaylar.tarihsayı} tarih ve sayılı dilekçesi"` |
| `9` | `Ceza Hafifleten Nedenler` | `hafifletennedenlar` | Free text (e.g. `"İlk defa disiplin suçu işlemesi, pişman olması"`) |
| `10` | `Ceza Şiddetlendiren Nedenler` | `siddetlendirennedenler` | Auto-filled by `AldıgiCezalar(adsoyad)` |
| `11` | `Kurulun Kanaati` | `kurulkanaati` | ComboBox: `"SUÇLU OLDUĞU YÖNÜNDE"` (default), `"SUÇSUZ OLDUĞU YÖNÜNDE"` |
| `12` | `Yönetmelik Maddesi` | `yönmaddesi` | ComboBox of `cezamaddeleri.madde` where `okultürü='LİSE'`. Selecting a clause auto-fills Rows 13 & 14. |
| `13` | `Maddeye Göre Verilecek Ceza` | `yönceza` | ReadOnly; auto-filled with `cezamaddeleri.ceza` |
| `14` | `Yönetmelik Maddesi Açıklaması` | `yönaçıklama` | ReadOnly; auto-filled with `cezamaddeleri.açıklama` |
| `15` | `Kurulun Verdiği Ceza` | `verilenceza` | ComboBox: `"SÖZLÜ OLARAK UYARMA"`, `"YAZILI OLARAK UYARMA"`, `"TOPLUM HİZMETİ"`, `"KINAMA"`, `"OKULDAN KISA SÜRELİ UZAKLAŞTIRMA"`, `"OKUL DEĞİŞTİRME"`, `"ÖRGÜN EĞİTİM DIŞINA ÇIKARMA"` |
| `16` | `Uzaklaştırma Cezası Gün Sayısı` | `cezagünsayısı` | Enabled ONLY when Row 15 == `"OKULDAN KISA SÜRELİ UZAKLAŞTIRMA"`. ComboBox: `"1 GÜN"`, `"2 GÜN"`, `"3 GÜN"`, `"4 GÜN"`, `"5 GÜN"`. Cleared (`""`) for all other penalties. |
| `17` | `Uzaklaştırma Cezası Uygulama Tarihi` | `cezauygulamatarihi` | Enabled ONLY when Row 15 == `"OKULDAN KISA SÜRELİ UZAKLAŞTIRMA"`. Cleared (`""`) otherwise. |
| `18` | `İndirilen Davranış Puanı` | `inddavranışpuanı` | ReadOnly; auto-calculated from Row 15 (see Section 5.6) |
| `19` | `Kalan Davranış Puanı` | `davranışpuanı` | ReadOnly; auto-calculated as `max(0, DavPuanı - Row18)` |
| `20` | `Kurul Kararı` | `karar` | Auto-generated by `KurulKararı()` whenever Rows 9–17 change, or manually editable |

### 5.6 LİSE Behavior Point Calculation (`Frm07KararLise2.cs` lines 1630–1695)
1. When a student is selected, their base behavior score `DavPuanı` is loaded (starts at `100`).
2. Whenever Row 15 (`verilenceza`) changes:
   - `"KINAMA"` -> `inddavranışpuanı = 10`
   - `"OKULDAN KISA SÜRELİ UZAKLAŞTIRMA"` -> `inddavranışpuanı = 20`
   - `"OKUL DEĞİŞTİRME"` -> `inddavranışpuanı = 40`
   - `"ÖRGÜN EĞİTİM DIŞINA ÇIKARMA"` -> `inddavranışpuanı = 80`
   - `"SÖZLÜ OLARAK UYARMA"`, `"YAZILI OLARAK UYARMA"`, `"TOPLUM HİZMETİ"`, or `""` -> `inddavranışpuanı = 0`
3. Remaining score:
   ```text
   Kalan = DavPuanı - inddavranışpuanı
   if Kalan < 0 then Kalan = 0
   ```
4. When saved (`BtnKaydet_Click`), `davranışpuanı` is written to both `OlayÖgrenciler.davranışpuanı` and `Ögrenciler.davranışpuanı`.

### 5.7 LİSE Auto-Generated Verdict Text (`Frm07KararLise2.KurulKararı()`)
- **Case 1 — Not Guilty (`Row[11] != "SUÇLU OLDUĞU YÖNÜNDE"`)**:
  - Clears Rows 12, 13, 14, 15, 16, 17, 18.
  - Sets `Row[20]` (`karar`) to:
    ```text
    {KURULADI} KURULUMUZCA YAPILAN İNCELEME, ARAŞTIRMA VE SORUŞTURMA NETİCESİNDE SUÇSUZ OLDUĞUNA KANAAT GETİRİLMİŞ, HERHANGİ BİR CEZAİ İŞLEM UYGULANMAMASINA OY BİRLİĞİ İLE KARAR VERİLMİŞTİR.
    ```
- **Case 2 — Guilty (`Row[11] == "SUÇLU OLDUĞU YÖNÜNDE"`)**:
  - **Subcase 2a — `Row[15] == "SÖZLÜ OLARAK UYARMA"`**:
    ```text
    MEB ORTAÖĞRETİM KURUMLARI YÖNETMELİĞİNİN {Row[12]} ({Row[14]}) BENDİNE GÖRE {Row[13]} CEZASINI GEREKTİREN FİİL VE DAVRANIŞLARDAN OLDUĞUNDAN  ANCAK {Row[9]} {Row[10]} NEDENLERİ DİKKATE ALINARAK SÖZLÜ OLARAK UYARILMASINA OY BİRLİĞİ İLE KARAR VERİLMİŞTİR.
    ```
  - **Subcase 2b — `Row[15] == "TOPLUM HİZMETİ"`**:
    ```text
    MEB ORTAÖĞRETİM KURUMLARI YÖNETMELİĞİNİN {Row[12]} ({Row[14]}) BENDİNE GÖRE {Row[13]} CEZASINI GEREKTİREN FİİL VE DAVRANIŞLARDAN OLDUĞUNDAN  ANCAK {Row[9]} {Row[10]} NEDENLERİ DİKKATE ALINARAK BİR TOPLUM HİZMETİ VERİLMESİNE OY BİRLİĞİ İLE KARAR VERİLMİŞTİR.
    ```
  - **Subcase 2c — Actual penalty equals clause's default penalty (`Row[13] == Row[15]`)**:
    - Line 1:
      ```text
      1. MEB ORTAÖĞRETİM KURUMLARI YÖNETMELİĞİNİN {Row[12]} ({Row[14]}) BENDİNE GÖRE {Row[13]} CEZASINI GEREKTİREN FİİL VE DAVRANIŞLARDAN OLDUĞUNDAN {Row[16]} {Row[15]} CEZASI VERİLMESİNE,
      ```
    - If `Row[16] != ""` (i.e. suspension days specified), appends Line 2:
      ```text
      2. CEZANIN {Row[17]} TARİHİNDE/TARİHLERİNDE UYGULANMASINA
      ```
    - Final line:
      ```text
      OY BİRLİĞİ İLE KARAR VERİLMİŞTİR.
      ```
  - **Subcase 2d — Actual penalty differs from clause's default penalty (`Row[13] != Row[15]`, e.g. mitigated to a lower penalty)**:
    - Line 1:
      ```text
      1. MEB ORTAÖĞRETİM KURUMLARI YÖNETMELİĞİNİN {Row[12]} ({Row[14]}) BENDİNE GÖRE {Row[13]} CEZASINI GEREKTİREN FİİL VE DAVRANIŞLARDAN OLDUĞUNDAN  ANCAK {Row[9]} {Row[10]} NEDENLERİ DİKKATE ALINARAK {Row[16]} {Row[15]} CEZASI VERİLMESİNE,
      ```
    - If `Row[16] != ""` (suspension days specified), appends Line 2:
      ```text
      2. CEZANIN {Row[17]} TARİHİNDE/TARİHLERİNDE UYGULANMASINA,
      ```
    - Final line:
      ```text
      OY BİRLİĞİ İLE KARAR VERİLMİŞTİR.
      ```

### 5.8 ORTAOKUL Verdict Grid & Text Generation (`Frm07KararOrta2.cs` — 17 Rows, 0-indexed `0..16`)
- Rows `0..14` match LİSE (with Ortaokul terminology: `"Davranış Yeri"`, `"Davranış Tarihi ve Saati"`, `"Davranış Çeşidi"`, `"Davranışın Nedeni"`, `"Yaptırım Hafifleten Nedenler"`, `"Yaptırım Şiddetlendiren Nedenler"`, `"Kurulun Kanaati"`: `["DAVRANIŞIN HATALI OLDUĞU YÖNÜNDE", "DAVRANIŞIN NORMAL OLDUĞU YÖNÜNDE"]`).
- Row `15` (`Kurulun Verdiği Yaptırım` / `verilenceza`): ComboBox `["UYARMA", "KINAMA", "OKUL DEĞİŞTİRME"]`.
- Row `16` (`Kurul Kararı` / `karar`):
  - **If `Row[11] != "DAVRANIŞIN HATALI OLDUĞU YÖNÜNDE"`**:
    ```text
    ÖĞRENCİ DAVRANIŞLARINI DEĞERLENDİRME KURULUMUZCA YAPILAN İNCELEME, ARAŞTIRMA VE SORUŞTURMA NETİCESİNDE ÖĞRENCİNİN DAVRANIŞININ HATALI OLMADIĞINA KANAAT GETİRİLMİŞ, HERHANGİ BİR YAPTIRIM UYGULANMAMASINA OY BİRLİĞİ İLE KARAR VERİLMİŞTİR.
    ```
  - **If `Row[11] == "DAVRANIŞIN HATALI OLDUĞU YÖNÜNDE"` and `Row[15]` is in `["UYARMA", "KINAMA", "OKUL DEĞİŞTİRME"]`**:
    ```text
    MEB OKUL ÖNCESİ EĞİTİM VE İLKÖĞRETİM KURUMLARI YÖNETMELİĞİNİN {Row[12]} ({Row[13]}) BENDİNE GÖRE {Row[15]} YAPTIRIMINI GEREKTİREN FİİL VE DAVRANIŞLARDAN OLDUĞUNDAN  ANCAK {Row[9]} {Row[10]} NEDENLERİ DİKKATE ALINARAK {Row[15]} YAPTIRIMI UYGULANMASINA OY BİRLİĞİ İLE KARAR VERİLMİŞTİR.
    ```
    *(Historical implementation note: In `Frm07KararOrta2.cs` line 1361, the original code interpolated `Rows[13]` (`yönceza`) inside the parentheses instead of `Rows[14]` (`yönaçıklama`).)*

---

## 6. Exact iTextSharp PDF Generation Specifications

### 6.1 Global PDF Conventions
1. **Library**: `iTextSharp.text` & `iTextSharp.text.pdf` (v5.x).
2. **Fonts**:
   - Primary font: `BaseFont.CreateFont("C:\\windows\\fonts\\tahoma.ttf", "Identity-H", embedded: true)`
   - Monospace font (used in `Frm19SeneBasiSonuKarar` for aligned tables): `BaseFont.CreateFont("C:\\windows\\fonts\\cour.ttf", "Identity-H", embedded: true)`
   - Standard font styles: `1` = Bold (`Font.BOLD`), `0` = Normal (`Font.NORMAL`), `2` = Italic (`Font.ITALIC`).
3. **Page Size & Margins**:
   - Portrait A4: `new Document(PageSize.A4, 30f, 20f, 20f, 20f)` (`Left=30, Right=20, Top=20, Bottom=20`), or `20f, 20f, 20f, 20f` on list reports.
   - Landscape A4: `new Document(PageSize.A4.Rotate(), 20f, 20f, 20f, 20f)`.
4. **Alignment Constants (`PdfPCell`)**:
   - `HorizontalAlignment`: `0` = Left (`ALIGN_LEFT`), `1` = Center (`ALIGN_CENTER`), `2` = Right (`ALIGN_RIGHT`), `3` = Justified (`ALIGN_JUSTIFIED`).
   - `VerticalAlignment`: `4` = Top (`ALIGN_TOP`), `5` = Middle (`ALIGN_MIDDLE`), `6` = Bottom (`ALIGN_BOTTOM`).
5. **File Naming & Lock Handling Loop**:
   - Every PDF generator starts at `num = 0` and builds `{Application.StartupPath}\{baseName}{num}.pdf`.
   - If the file exists, it attempts `Process.Start(path); process.Kill(); File.Delete(path);`. If locked by a PDF viewer (throws `Exception`), it increments `num++` and retries until an available filename is found.

---

### 6.2 `Frm05ifadeTutanaklari.cs` — 14 Investigation & Counseling Forms
All 14 forms use Portrait A4 (`30f, 20f, 20f, 20f`) and Tahoma `Identity-H`:
- `font` = 12pt Bold, `font2` = 9pt Bold, `font3` = 8pt Bold, `font4` = 7pt Bold, `font5` = 8pt Normal, `font6` = 7pt Normal.

#### 1. `ifade-tutanagi{N}.pdf` (`Yazifade`) & 2. `savunma-tutanagi{N}.pdf` (`YazSavunma`)
- **Table 1 (5 columns, widths `[1f, 1f, 1f, 1f, 1f]`, `TotalWidth = PageWidth - Left - Right`)**:
  - **Title Cell (`Colspan=5`, Center, `MinimumHeight=22f`, `BorderWidth=0f`, 12pt Bold)**:
    - `"İFADE TUTANAĞI"` (in `Yazifade`) or `"SAVUNMA TUTANAĞI"` (in `YazSavunma`).
  - **Student Info Block (Left label `Colspan=1` 8pt Bold, Right value `Colspan=4` 8pt Normal, `BorderWidth=0f`)**:
    - `"İFADESİ ALINANIN"` / `"SAVUNMASI ALINANIN"` (`Colspan=5`, 8pt Bold)
    - `"Adı Soyadı"` : `": {AdSoyad}"`
    - `"Sınıfı ve Okul No"` : `": {Sınıf}  -  {OkulNo}"`
    - `"Baba - Anne Adı"` : `": {BabaAdı} - {AnneAdı}"`
    - `"Doğum Yeri ve Tarihi"` : `": {DoğumYeri} - {DoğumTarihi}"`
    - `"İfade Tarihi"` / `"Savunma Tarihi"` : `": {Now.ToShortDateString()}"`
  - **Event Prompt Cell (`Colspan=5`, `BorderWidth=0f`, `BorderWidthTop=0.1f`, `BorderWidthBottom=0.1f`, 8pt Bold)**:
    - `"OLAY: \r\n" + ifadeMetni`
  - **Ruled Writing Lines (`35` rows)**:
    - Loop `i = 1..35`: adds a cell with `Colspan=5`, `MinimumHeight=16f`, `BorderWidth=0f`, `BorderWidthBottom=0.1f` (thin horizontal ruled lines for handwritten statement, or pre-filled text if any).
  - **Student Signature Row**:
    - `"İFADEYİ VERENİN ADI SOYADI - İMZASI"` (or `"SAVUNMAYI VERENİN..."`), right-aligned, `Colspan=5`, `MinimumHeight=30f`.
  - **Board Member Signatures Row (5 cells, 1 per column, Center, `MinimumHeight=45f`, `BorderWidth=0f`, 7pt Normal)**:
    - Col 1: `"İfadeyi Alan\r\n{KurulBaşkanı}\r\n{KurulAdı} Kurulu Başkanı"`
    - Col 2..5: `"İfadeyi Alan\r\n{Üye1..4}\r\nÜye"`

#### 3. `bilgitoplamaolay-tutanagi{N}.pdf` (`YazBilgiToplamaOlay`)
- Header (`Colspan=5`, Center, 12pt Bold): `"{KURULADI} KURULU BİLGİ TOPLAMA FORMU - Olay"`
- Section 1 (`Colspan=5`, 8pt Bold): `"Olayla İlgili Bilgilerim:"` followed by **16 ruled rows** (`MinimumHeight=16f`, `BorderWidthBottom=0.1f`).
- Section 2 (`Colspan=5`, 8pt Bold): `"Arkadaşım Hakkındaki Görüşlerim:"` followed by **16 ruled rows** (`MinimumHeight=16f`, `BorderWidthBottom=0.1f`).
- Bottom right signature block: `"Öğrencinin Adı Soyadı / Sınıfı - Numarası / İmzası"`, `"Tarih: ..../..../20...."`.

#### 4. `bilgitoplamaögretmen-tutanagi{N}.pdf` (`YazBilgiToplamaÖgretmen`)
- Header: `"{KURULADI} KURULU BİLGİ TOPLAMA FORMU - Öğretmen"`
- Student identification rows (`Adı Soyadı`, `Sınıfı - Numarası`).
- Three sections, each with a bold heading and **12 ruled rows** (`MinimumHeight=16f`, `BorderWidthBottom=0.1f`):
  1. `"Dersimdeki Tutumu:"`
  2. `"Ders Dışı Faaliyetleri:"`
  3. `"Arkadaşları İle İlişkisi:"`
- Bottom right signature block: `"Öğretmenin Adı Soyadı / Branşı / İmzası"`.

#### 5. `bilgitoplamasınıfögretmen-tutanagi{N}.pdf` (`YazBilgiToplamaSınıfÖgretmen`)
- Header: `"{KURULADI} KURULU BİLGİ TOPLAMA FORMU - Sınıf Rehber Öğretmeni"`
- Student identification rows.
- Two sections, each with **18 ruled rows** (`MinimumHeight=16f`, `BorderWidthBottom=0.1f`):
  1. `"Kişisel Özellikleri:"`
  2. `"Sosyal Özellikleri:"`
- Bottom right signature block: `"Sınıf Rehber Öğretmeni Adı Soyadı / İmzası"`.

#### 6. `bilgitoplamaveli-tutanagi{N}.pdf` (`YazBilgiToplamaVeli`)
- Header: `"{KURULADI} KURULU VELİ BİLGİ TOPLAMA FORMU"`
- Student & Guardian identification rows.
- Heading: `"Öğrenci Hakkındaki Görüşlerim:"` followed by **34 ruled rows** (`MinimumHeight=16f`, `BorderWidthBottom=0.1f`).
- Bottom right signature block: `"Velinin Adı Soyadı / İmzası"`.

#### 7. `bilgitoplamaörenci-tutanagi{N}.pdf` (`YazBilgiToplamaÖgrenci`)
- Header: `"{KURULADI} KURULU BİLGİ TOPLAMA FORMU - Öğrenci"`
- Three sections, each with **12 ruled rows** (`MinimumHeight=16f`, `BorderWidthBottom=0.1f`):
  1. `"Dersteki Davranışları:"`
  2. `"Teneffüslerdeki Davranışları:"`
  3. `"Okul Dışındaki Davranışları:"`
- Bottom right signature block: `"Öğrencinin Adı Soyadı / İmzası"`.

#### 8. `eksüretalep-form{N}.pdf` (`YazEksüreTalep`)
- Header (`Colspan=2`, Center, 12pt Bold): `"{KURULADI} KURULU\r\nOLAY EK SÜRE TALEP FORMU"`
- **2-Column Bordered Table (`widths = [1f, 3f]`, `BorderWidth = 0.5f`, `MinimumHeight = 25f`)**:
  - `"Olay"` : `Olaylar.olayadı`
  - `"Olay Açıklaması"` : `Olaylar.özet`
  - `"Öğrenci Adı Soyadı"` : `{Sınıf} - {OkulNo} - {AdSoyad}`
  - `"Olaydaki Durumu"` : `OlayÖgrenciler.durumu`
  - `"Ek Süre"` : `"{TxtSüre0.Text} GÜN"`
  - `"Ek Süre Nedeni"` : `TxtSüre1.Text`
  - `"Tarih"` : `Now.ToShortDateString()`
- **Petition Block (`SpacingBefore=30f`)**:
  - Centered Bold Header: `"{OKULADI} MÜDÜRLÜĞÜNE"`
  - Body text (LİSE references *MEB Ortaöğretim Kurumları Yönetmeliğinin 188. maddesinin 3. fıkrası*).
  - Right-aligned signature: `"{KurulBaşkanı}\r\n{KurulAdı} Kurulu Başkanı"`
  - Centered bottom approval (`Colspan=2`): `"O L U R\r\n{Now.ToShortDateString()}\r\n\r\n\r\n{Müdüradı}\r\nOkul Müdürü"`.

#### 9. `sikayet-form{N}.pdf` (`YazSikayetFormu`)
- Header: `dsp_Modül.Antet` followed by `"Ş İ K A Y E T   F O R M U"` (12pt Bold).
- 2-Column Bordered Info Table (`widths = [2f, 5f]`, `BorderWidth=0.5f`, `MinimumHeight=22f`):
  - `"Şikayet Eden (Adı-Soyadı-Sınıfı-No)"`, `"Şikayet Edilen (Adı-Soyadı-Sınıfı-No)"`, `"Tanıklar"`, `"Deliller"`, `"Olay Tarihi ve Saati"`, `"Olayın Konusu"`.
- `"Olayın Açıklaması:"` followed by 5 dotted lines (`"........................................................................................................................................"`).
- Referral Box signed by `"{Müdüradı}\r\nOkul Müdürü"`.
- `"Rehberlik ve Psikolojik Danışma Servisi Görüş ve Önerisi:"` followed by 13 dotted lines and `"Rehber Öğretmen"` signature.

#### 10. `sözlü-uyari-tutanagi{N}.pdf` (`YazSözlüUyariTutanagi`)
- Header: `dsp_Modül.Antet` + `"S Ö Z L Ü   U Y A R I   T U T A N A Ğ I"`.
- Student Info Table (`widths = [1f, 3f]`, `Adı Soyadı`, `Sınıfı / Numarası`, `Tarih`, `Konu`).
- Body statement + 2-column signature table (`"UYARILAN ÖĞRENCİ"` / `"UYARAN ÖĞRETMEN / YÖNETİCİ"`).

#### 11. `ögrenci-sözlesmesi{N}.pdf` (`YazSözlesme`)
- Header: `dsp_Modül.Antet` + `"Ö Ğ R E N C İ    S Ö Z L E Ş M E S İ"`.
- Student Info Table + Pledge paragraph + Student signature block.

#### 12. `veli-görüsme-tutanagi{N}.pdf` (`YazVeliGörüsmeTutanagi`) & 13. `veli-görüsme-form{N}.pdf` (`YazVeliGörüsmeFormu`)
- Parent-teacher interview record and blank form with student/parent info tables, dotted lines for interview notes and decisions, and dual signature blocks (`"VELİ"` / `"GÖRÜŞMEYİ YAPAN"`).

#### 14. `rehber-ögretmen-görüs-form{N}.pdf` (`YazRehberÖgretmenGörüsFormu`)
- Header: `"{OKULADI} MÜDÜRLÜĞÜNE\r\n({KurulAdı} Kurulu Başkanlığına)"`
- Intro paragraph identifying the student (`{Sınıf} sınıfı {OkulNo} numaralı öğrencisi {AdSoyad}`).
- 4 numbered counseling sections (`1- Öğrencinin Ailevi Durumu:`, `2- Öğrencinin Psikolojik ve Sosyal Özellikleri:`, `3- Olayın Nedenleri Hakkındaki Görüşler:`, `4- Öğrenciye Uygulanabilecek Yaptırım / Ceza veya Rehberlik Önlemleri Hakkındaki Öneriler:`), each with 7 dotted lines.
- Signed at bottom right by `"Psikolojik Danışman / Rehber Öğretmen"`.

---

### 6.3 `Frm06CagriPusulasi.cs` — Summons Slips & Meeting Notice PDFs

#### 1. `ögrencisavunmacagripusulasi{N}.pdf` (`YazÖgrenciSavunmaCagri`) & 2. `ögrenciifadecagripusulasi{N}.pdf` (`YazÖgrenciİfadeCagri`)
- Generates **1 page per checked student** (`document.NewPage()` at the end of each student loop):
  - **Table 1 (5 columns `[1,1,1,1,1]`)**:
    - `dsp_Modül.Antet` (8pt Bold, Center, `Colspan=5`, `MinimumHeight=35f`)
    - `"{KURULADI} KURULU"` (9pt Bold, Center, `Colspan=5`, `MinimumHeight=22f`)
    - `"S A V U N M A     Ç A Ğ R I     P U S U L A S I"` or `"İ F A D E     Ç A Ğ R I     P U S U L A S I"` (20pt Bold, Center, `Colspan=5`, `MinimumHeight=22f`)
  - **Table 2 (2 columns `[1f, 3f]`, indented by 50f left/right, top & bottom border `2f`)**:
    - `"ÖĞRETİM YILI / DÖNEM:"` : `dsp_Modül.ÖgretimYılı`
    - `"SAVUNMANIN ALINACAĞI YER :"` / `"İFADENİN ALINACAĞI YER :"` : `TxtYer.Text`
    - `"TARİH VE SAATİ :"` : `"{DTtarih.Value.Date.ToLongDateString()}  - {TxtSaat.Text}"`
  - **Table 3 (2 columns `[1f, 3f]`, indented by 50f left/right)**:
    - `"Öğrencinin;"`
    - `"ADI SOYADI : "` : `ÖgrenciAd`
    - `"SINIFI - NUMARASI :"` : `ÖgrenciSınıfNo`
  - **Table 4 (3 columns `[1f, 1f, 1f]`, indented by 20f left/right)**:
    - Summons paragraph (`Colspan=3`, 8pt Normal — LİSE cites *Madde 190*; ORTAOKUL cites *ilgili madde*).
    - Signature cell (`Colspan=1` on right after `Colspan=2` empty cell, `MinimumHeight=77f`, Center):
      `"{KurulBaşkanı}\r\n{KurulAdı} Kurulu Başkanı\r\n{KurulBaşkanıünvan}"`
  - **Table 5 (2 columns `[1f, 3f]`, `"OKUDUM"` acknowledgement block)**:
    - `"OKUDUM"`, `"ADI SOYADI:" : ".................................."`, `"TARİH:" : "....../....../20........."`, `"İMZA:" : ".................................."`.
  - **Table 6 (Footer positioned via `WriteSelectedRows` at `BottomMargin + 30f`, `BorderWidthTop = 2f`, 7pt Italic)**:
    - `"ADRES: {Adres}\r\nTEL NO: {TelNo}\r\nFAX NO: {FaxNo}"`

#### 3. `üyecagripusulasi{N}.pdf` (`YazÜyeCagri`)
- **Table 1**: `dsp_Modül.Antet` (Centered).
- **Table 2 (3 columns `[1f, 14f, 2f]`)**:
  - `"Sayı:"`, `TxtSayı.Text`, `TxtTarih.Value.Date` (Right-aligned)
  - `"Konu:"`, `TxtKonu.Text`, `""`
- **Table 3 (3 columns `[1f, 1f, 1f]`)**:
  - Header (`Colspan=3`, Center, 8pt Bold): `"{KURULADI} KURULU ÜYELERİNE"`
  - Body (`Colspan=3`, 8pt Normal): Meeting invitation paragraph with `DTtarih`, `TxtSaat`, `TxtYer`.
  - Signatures: Left cell `"{KurulBaşkanı}\r\n{KurulAdı} Kurulu Başkanı\r\n{KurulBaşkanıünvan}"`, Right cell `"{Müdüradı}\r\nOkul Müdürü"`.
- **Table 4 (1 column, indented by 20f)**:
  - `"G Ü N D E M"` (8pt Bold) followed by numbered agenda rows `"{sıra}) {madde}"` from `GridGündem`.
- **Table 5 (2 columns `[3f, 1f]`, indented by 100f)**:
  - Header: `"ÜYELER"` | `"İMZA"`
  - 4 rows for `Üye1`..`Üye4` with `".................................."` signature lines.
- **Table 6**: Absolute-positioned bottom address/phone/fax footer (`BorderWidthTop = 2f`).

#### 4. `cagri-pusulasi-liste{N}.pdf` (`YazdirListe`)
- Title: `"{OkulAdı}\r\nÇAĞRI PUSULASI ÖĞRENCİ LİSTESİ"` (12pt Bold, Center).
- 3-Column Table (`widths = [1f, 5f, 4f]`, `BorderWidth = 0.3f`, header height `30f`, row height `15f`):
  - Columns: `SIRA`, `ADI SOYADI`, `SINIFI - NUMARASI` for all checked students, signed by `KurulBaşkanı` at bottom right.

---

### 6.4 `Frm06YaziliUyari.cs` — `yazili-uyari{N}.pdf` (`YazdirUyari`)
- Portrait A4 (`30f, 20f, 20f, 20f`).
- **Table 1**: `dsp_Modül.Antet` (Center, 8pt Bold).
- **Table 2 (`widths = [1f, 14f, 3f]`)**: `"Sayı:"` (`TxtSayi.Text`), `DTTarih.Value.ToShortDateString()`, `"Konu: Yazılı Uyarı"`.
- **Table 3**:
  - Title: `"ÖĞRENCİ YAZILI UYARI TUTANAĞI"` (12pt Bold, Center).
  - Student details (`Adı Soyadı`, `Sınıfı - Okul No`, `Uyarı Tarihi`).
  - Warning body citing *MEB Ortaöğretim Kurumları Yönetmeliği Madde {Madde} ({Açıklama})* and *Madde 163/2*.
  - Right-aligned signature: `"{Müdüradı}\r\nOkul Müdürü"`.
  - Student acknowledgement section (`"Yukarıdaki yazılı uyarıyı okudum ve bir nüshasını elden teslim aldım..."`) with `"Öğrencinin Adı Soyadı / İmzası"` and `"Tebliğ Eden Müdür Yardımcısı"`.

---

### 6.5 `Frm07KararLise2.cs` & `Frm07KararOrta2.cs` — `karar-tutanak{N}.pdf` (`EK-1` & `EK-10`)

#### 1. LİSE Verdict PDF (`Frm07KararLise2.YazKararTutanak` — `EK-1`)
- Portrait A4 (`30f, 20f, 20f, 20f`), Tahoma `Identity-H`.
- **Dynamic Row Height Calculation**:
  - Sums the total text length of all rows in `GridÖgrenciBilgileri` and `GridKarar` to compute a vertical scaling multiplier `num6` (`MinimumHeight = 16 * num6`, or `12 * num6` if total characters > 1300) so the entire `EK-1` form fits onto a single A4 page without overflowing.
- **Header Table (`pdfPTable`, 3 columns `[12f, 1f, 20f]`)**:
  - Row 1 (`Colspan=3`, Right-aligned, 9pt Bold): `"EK-1"`
  - Row 2 (`Colspan=3`, Center, 12pt Bold): `"{KURULADI} KURULU KARARI"`
- **Main Body Table (`pdfPTable`, 3 columns `[12f, 1f, 20f]`, thin cell borders `0.1f`)**:
  - Renders rows from `GridÖgrenciBilgileri` (Karar No, Karar Tarihi, Öğrencinin Adı Soyadı, Sınıfı, Okul No, Cinsiyeti, TC No, Doğum Yeri/Tarihi, Baba/Anne Adı, Veli Adı/Tel/Adres, Yatılılık, Başarı, Sağlık, Ekonomik Durum, Nerede Kaldığı, Anne-Baba Sağ/Öz, Nerede Okuduğu, Büyüdüğü Çevre, Daha Önce Ceza Aldı mı) + `GridKarar` rows (Suç Yeri, Suç Tarihi, Suç Çeşidi, Suçun Nedeni, Olayın Özeti, Tanıkların İfadesi, Deliller, Hafifleten Nedenler, Şiddetlendiren Nedenler, Kurulun Kanaati, Yönetmelik Maddesi, Verilen Ceza, Uzaklaştırma Gün Sayısı & Tarihi, İndirilen/Kalan Davranış Puanı, and **Kurul Kararı**).
  - Left column (`width=12f`): Field label (8pt Bold).
  - Middle column (`width=1f`): `":"` (8pt Bold, Center).
  - Right column (`width=20f`): Field value (8pt Normal, or 8pt Bold for `Kurul Kararı`).
- **Signature Table (`pdfPTable2`, 5 columns `[1f, 1f, 1f, 1f, 1f]`, `SpacingBefore=10f`)**:
  - Row 1 (`MinimumHeight=25f`, Center, 7pt Bold): `"Başkan"` | `"Üye"` | `"Üye"` | `"Üye"` | `"Üye"`
  - Row 2 (`MinimumHeight=35f`, Center, 7pt Normal):
    - Col 1: `"{KurulBaşkanı}\r\n{KurulBaşkanıünvan}"`
    - Col 2..5: `"{Üye1..4}\r\n{Üye1..4ünvan}"`
  - Row 3 (`Colspan=5`, Center, `MinimumHeight=45f`, 8pt Bold):
    - `"....../....../20......\r\nUygundur\r\n\r\n{Müdüradı}\r\nOkul Müdürü"`

#### 2. ORTAOKUL Verdict PDF (`Frm07KararOrta2.YazKararTutanak` — `EK-10`)
- Same 3-column `[12f, 1f, 20f]` outer structure with `"EK-10"` top-right label and `"ÖĞRENCİ DAVRANIŞLARINI DEĞERLENDİRME KURULU KARARI"` title.
- Includes a **nested 3-column table (`widths = [12f, 10f, 10f]`)** for the Mother/Father comparison section:
  - Header: `""` | `"ANNE"` | `"BABA"`
  - Rows: `"Adı Soyadı"`, `"Yaşı"`, `"Eğitim Durumu"`, `"Hayatta mı?"`, `"Öz mü?"`, `"Birlikte mi?"`.
- Followed by the same 5-column Board Signature block + `"Okul Müdürü"` approval block.

---

### 6.6 `Frm08CezaTabligati.cs` — Student & Parent Penalty Notification PDFs

#### 1. `ögrencitabligat{N}.pdf` (`YazÖgrenciTebligatPDF`)
- Portrait A4 (`30f, 20f, 20f, 20f`). Generates 1 page per checked student:
  - **Table 1 (5 columns)**: `dsp_Modül.Antet` + `"{KURULADI} KURULU"` + `"Ö Ğ R E N C İ   K A R A R   T E B L İ G A T   P U S U L A S I"` (16pt Bold).
  - **Table 2 (2 columns `[1f, 2f]`, indented by 50f, top & bottom border `2f`)**:
    - `"ÖĞRETİM YILI"` : `dsp_Modül.ÖgretimYılı`
    - `"KARAR TARİHİ VE NO"` : `"{KararTarihi}  -  {KararNo}"`
    - `"ÖĞRENCİNİN ADI SOYADI"` : `AdSoyad`
    - `"SINIFI VE OKUL NO"` : `SınıfNo`
    - `"VERİLEN CEZA"` / `"UYGULANAN YAPTIRIM"` : `VerilenCeza`
    - `"CEZA UYGULAMA TARİHİ"` (if suspension) : `CezaUygulamaTarihi`
    - `"CEZA NEDENİ DAVRANIŞ"` : `CezaNedeni`
    - `"YÖNETMELİK MADDESİ"` : `"{YönMadde} ({YönAçıklama})"`
  - **Table 3 (3 columns `[1f, 1f, 1f]`)**:
    - Notice paragraph (LİSE cites *Madde 169*; ORTAOKUL cites *ilgili madde*).
    - Right-aligned signature: `"{Müdüradı}\r\nOkul Müdürü"`.
  - **Table 4 (2 columns `[1f, 1f]`, `"TEBLİĞ EDEN"` / `"TEBELLÜĞ EDEN (ÖĞRENCİ)"`)**:
    - Prompts for `Adı Soyadı`, `Tarih`, `İmza`.
  - **Table 5**: Bottom address/phone/fax footer (`BorderWidthTop = 2f`).

#### 2. `Velitabligat{N}.pdf` (`YazVeliTebligat`)
- Portrait A4 (`30f, 20f, 20f, 20f`). Generates 1 page per checked student:
  - **Table 1**: `dsp_Modül.Antet`.
  - **Table 2 (`widths = [1f, 14f, 3f]`)**:
    - `"SAYI:"` : `"{ResmiYazıKodu}-{sequential_number}"` | `Now.ToShortDateString()`
    - `"KONU:"` : `"Ceza Bildirimi"`
  - **Table 3**:
    - Salutation: `"SAYIN; {VeliAdı}"` (8pt Bold)
    - Parent notification text (cites 5 business days right of objection under *Madde 170 ve 160* for LİSE or *Madde 56/10* for ORTAOKUL).
    - Right-aligned signature: `"{Müdüradı}\r\nOkul Müdürü"`.
  - **Table 4 (4 columns `[2f, 2f, 3f, 5f]`, `BorderWidth = 0.5f`)**:
    - Headers: `"KARAR TARİH VE NO"` | `"VERİLEN CEZA"` | `"YÖNETMELİK MADDESİ"` | `"CEZA NEDENİ DAVRANIŞ"`
    - 1 data row with the student's verdict details.
  - **Table 5 (`"TEBLİĞ EDEN"` / `"TEBELLÜĞ EDEN (VELİ)"`)** + **Table 6 (Footer)**.

---

### 6.7 `Frm10DavranisPuaniCizelge.cs` & `Frm10yaptırımSilme.cs` — Landscape End-of-Year Reports

#### 1. `davranispuaniiadecizelgesi{N}.pdf` (`Frm10DavranisPuaniCizelge`)
- **Orientation**: Landscape A4 (`PageSize.A4.Rotate(), 20f, 20f, 20f, 20f`).
- **Header (`Colspan=14`, Center, 11pt Bold)**:
  - `"{OkulAdı}\r\n{ÖgretimYılı} EĞİTİM ÖĞRETİM YILI\r\nMEB ORTAÖĞRETİM KURUMLARI YÖNETMELİĞİNİN 171/2 MADDESİ GEREĞİ SENE SONU DAVRANIŞ PUANI İADE ÇİZELGESİ"`
- **Table (14 columns, widths `[2f, 6f, 3f, 4f, 5f, 3f, 3f, 6f, 3f, 3f, 3f, 3f, 4f, 4f]`, `BorderWidth = 0.3f`)**:
  - Columns: `SIRA`, `ADI SOYADI`, `OKUL NO`, `SINIFI`, `OLAY ADI`, `KARAR NO`, `KARAR TARİHİ`, `VERİLEN CEZA`, `YÖN. MADDESİ`, `İNDİRİLEN PUAN`, `KALAN PUAN`, `İADE PUAN`, `İADE TARİHİ`, `İADE DURUMU`.
- **Footer Signatures**: 5-column Board Signature table (`KurulBaşkanı`, `Üye1..4`) + centered `Okul Müdürü` approval.

#### 2. `yaptirim-karari-kaldirma-raporu{N}.pdf` (`Frm10yaptırımSilme`)
- **Orientation**: Landscape A4 (`PageSize.A4.Rotate(), 20f, 20f, 20f, 20f`).
- **Header**: `"{OkulAdı}\r\n{ÖgretimYılı} EĞİTİM ÖĞRETİM YILI\r\nMEB OKUL ÖNCESİ EĞİTİM VE İLKÖĞRETİM KURUMLARI YÖNETMELİĞİNİN 62/3 MADDESİ GEREĞİ YAPTIRIM KARARI KALDIRMA RAPORU"`
- **Table (11 columns, widths `[2f, 6f, 3f, 4f, 6f, 3f, 4f, 6f, 4f, 4f, 5f]`, `BorderWidth = 0.3f`)**:
  - Columns: `SIRA`, `ADI SOYADI`, `OKUL NO`, `SINIFI`, `OLAY ADI`, `KARAR NO`, `KARAR TARİHİ`, `UYGULANAN YAPTIRIM`, `YÖN. MADDESİ`, `KALDIRMA TARİHİ`, `DURUMU`.
- **Footer Signatures**: 5-column Board Signature table + centered `Okul Müdürü` approval.

---

### 6.8 `Frm11Cezaistatistik.cs`, `Frm12CezaAlanlar.cs`, `Frm14Dizipusulasi.cs`, `Frm19SeneBasiSonuKarar.cs`

#### 1. `cezaistatistik{N}.pdf`, `olayistatistik{N}.pdf`, `yönetmelikistatistik{N}.pdf` (`Frm11Cezaistatistik`)
- Portrait or Landscape A4 depending on `CBoxSayfaYönü`.
- If `CBoxGrafik == "GRAFİK EKLE"`, saves the MSChart control to `{StartupPath}\grafik.png` and embeds it into a 4th table column (`widths = [1f, 6f, 1f, 8f]`) with `Rowspan = grid.Rows.Count`; otherwise uses a 3-column table (`widths = [1f, 6f, 1f]`: `SIRA`, `CEZA`/`OLAY`/`MADDE`, `SAYI`).
- Followed by 5-column Board Signature block + `Okul Müdürü`.

#### 2. `sınıfcezaalanögrencilistesi{N}.pdf` (`Frm12CezaAlanlar`)
- Portrait A4 (`20f, 20f, 20f, 20f`).
- Title: `"{OkulAdı}\r\nSINIF CEZA ALAN ÖĞRENCİLER RAPORU"` (Tab 1) or `"VERİLEN CEZAYA GÖRE ÖĞRENCİ LİSTESİ"` (Tab 2).
- 8-Column Table (`widths = [1f, 2f, 2f, 3f, 3f, 6f, 2f, 2f]`, `BorderWidth = 0.3f`):
  - `SIRA`, `OKUL NO`, `SINIFI`, `ADI SOYADI`, `OLAY ADI`, `VERİLEN CEZA`, `KARAR NO`, `KARAR TARİHİ`.

#### 3. `{N}.pdf` — Dizi Pusulası (`Frm14Dizipusulasi`)
- Portrait A4 (`30f, 20f, 20f, 20f`).
- Title: `"{OKULADI} MÜDÜRLÜĞÜ\r\n{KURULADI} KURULU DİZİ PUSULASI"` (12pt Bold, Center).
- 3-Column Table (`widths = [1f, 6f, 2f]`, indented by 80f, `BorderWidth = 0.3f`, row height `22f`):
  - Headers: `SIRA` | `KONUSU` | `ADET / SAYI`
  - Appends a `"TOPLAM"` row summing `Conversion.Val(Cells[2].Value)`.
- Bottom signature line: `"Teslim Eden                                                                                          Teslim Alan"`.

#### 4. `sene-basi-sonu-tutanak{N}.pdf` (`Frm19SeneBasiSonuKarar.YazdirTutanak`)
- Portrait A4 (`30f, 20f, 20f, 20f`). Uses both `tahoma.ttf` and `cour.ttf` (Courier New 10pt Normal/Bold for decision numbers and body so fixed-width ASCII tables align).
- **Table 1 (1 column)**: `"{KURULADI} KURULU {CBOXtopAdi.Text} TOPLANTI TUTANAĞI"` (16pt Bold, Center).
- **Table 2 (2 columns `[1f, 4f]`, `BorderWidth = 0f`)**:
  - `"OKUL ADI"` : `": {OkulAdı}"`
  - `"ÖĞRETİM YILI"` : `": {ÖgretimYılı}"`
  - `"TOPLANTI TARİHİ"` : `": {DTPtoptarihi.Value.ToLongDateString()}"`
  - `"TOPLANTI SAATİ"` : `": {MTBOXtopsaati.Text}"`
  - `"TOPLANTI YERİ"` : `": {TBOXtopyeri.Text}"`
  - `"TOPLANTI BAŞKANI"` : `": {CBoxTopBaskani.Text}"`
- **Table 3 (1 column)**:
  - Statutory intro paragraph (LİSE cites *MEB Ortaöğretim Kurumları Yönetmeliğinin 189. Maddesi*; ORTAOKUL cites *MEB Okul Öncesi Eğitim ve İlköğretim Kurumları Yönetmeliğinin 58. Maddesi*).
- **Table 4 (2 columns `[1f, 14f]`)**:
  - For each non-empty row in `dgvKararlar`:
    - Col 1: `"{NO}.) "` (Courier New 10pt Bold, Right-aligned)
    - Col 2: `"{KARAR}\r\n"` (Courier New 10pt Normal, Left-aligned)
  - Closing row (`Colspan=2`, Tahoma 11pt Bold):
    - If `SENE BAŞI`: `"OY BİRLİĞİ İLE KARAR VERİLDİ."`
    - If `SENE SONU`: `"İŞ BU KURUL TUTANAĞININ SENE SONU ÖĞRETMENLER KURULUNDA OKUNMASI TALEBİYLE OKUL MÜDÜRLÜĞÜNÜN ONAYINA SUNULMASINA OY BİRLİĞİ İLE KARAR VERİLDİ."`
- **Table 5 (5 columns `[1f, 1f, 1f, 1f, 1f]`)**:
  - 5 Board Member signatures (`KurulBaşkanı`, `Üye1..4`) + centered `Okul Müdürü` approval block (`"....../....../20......\r\n\r\n{Müdüradı}\r\nOkul Müdürü"`).
