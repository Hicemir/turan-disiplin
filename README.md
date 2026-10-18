# 🏛️ Okul Evrak - Öğrenci Disiplin ve Davranış Değerlendirme Sistemi

MEB Ortaöğretim ve İlköğretim Kurumları Yönetmeliklerine tam uyumlu, modern, hızlı ve kullanımı son derece kolay masaüstü disiplin ve evrak yönetim sistemi.

---

## 📥 Programı İndir (Kurulumsuz Tek Dosya)

Program herhangi bir kurulum sihirbazı veya harici veritabanı kurulumu gerektirmez. İndirdikten sonra doğrudan çift tıklayarak çalıştırabilirsiniz:

| İşletim Sistemi | İndirme Linki | Çalıştırma |
| :--- | :--- | :--- |
| **🪟 Windows (10 / 11 / 8 / 7)** | [📥 **OkulEvrak.exe İndir**](executables/OkulEvrak.exe?raw=true) | İndirdiğiniz `OkulEvrak.exe` dosyasına çift tıklayarak hemen kullanmaya başlayabilirsiniz. |
| **🐧 Linux (Ubuntu, Debian, Fedora vb.)** | [📥 **OkulEvrak-Linux İndir**](executables/OkulEvrak-Linux?raw=true) | İndirdikten sonra çift tıklayın veya terminalden `./OkulEvrak-Linux` yazarak çalıştırın. |

*(GitHub Releases sayfasından da en son sürüme ait `.exe` veya Linux çalıştırılabilir dosyalarını tek tıkla indirebilirsiniz).*

---

## ✨ Öne Çıkan Özellikler

* **🎨 İki Farklı Arayüz Modu**:
  * **Modern Arayüz**: Büyük, okunabilir kartlar ve adım adım yönlendirme.
  * **Klasik Arayüz**: Orijinal .NET programının birebir 1:1 nostaljik WinForms masaüstü tasarımı.
* **🌈 Tema ve Yazı Boyutu Tercihleri**:
  * **Açık Tema**: Gündüz kullanımı için yüksek kontrastlı net görünüm.
  * **Sarı / Krem**: Orijinal programın sıcak nostaljik krem tonu.
  * **Koyu Tema**: Göz yormayan koyu gri/gece modu.
  * **Yazı Boyutları**: `15pt`, `17pt` ve `19pt` (kolay okunabilirlik).
* **👥 e-Okul & e-Mesem Excel Aktarımı**:
  * Standart e-Okul sınıf listeleri (.xlsx / .xls)
  * e-Mesem çıraklık/stajyer listeleri
  * 3 sütunlu basit hızlı listeler
* **⚖️ %100 Mevzuata Uygun Resmi Evrak ve Karar Hesaplama**:
  * Lise için **EK-1**, Ortaokul için **EK-10** Karar Tutanakları.
  * 100 üzerinden davranış puanı düşümü, hafifletici nedenler ve ceza gerekçelendirmeleri.
* **📄 Tek Tıkla A4 Resmi PDF Üretimi**:
  1. **Disiplin Kurulu Kararı (EK-1 / EK-10)** - 21/17 satırlı resmi format.
  2. **Çağrı Pusulaları** - İfade ve savunma çağrısı, tebellüğ bloğu ile.
  3. **Resmi İfade ve Savunma Tutanakları**.
  4. **Öğrenci Ceza Tebligat Pusulası** - 5 günlük yasal itiraz hakkı metniyle.
  5. **Veli Ceza Tebligat Pusulası**.
  6. **Sene Başı / Sene Sonu Kurul Toplantı Tutanağı**.
  7. **Soruşturma Dosyası Dizi Pusulası**.
  8. **Madde 171/2 Sene Sonu Puan İade Çizelgesi** / Yaptırım Kaldırma Raporu.
  9. **Dönemlik Ceza ve Olay İstatistik Raporları**.
  10. **Madde 163/2 Yazılı Uyarı Tutanağı**.
* **⚙️ Yönetmelik ve Şablon Düzenleyici**:
  * Yeni suç ve ceza maddeleri ekleyebilme, mevcut maddeleri özelleştirebilme.
  * `os.UserConfigDir()` üzerinde kullanıcıya özel güvenli saklama.

---

## 🚀 Geliştirici & Kaynak Koddan Derleme

Projeyi yerel makinenizde kaynak koddan derlemek isterseniz:

```bash
# Bağımlılıkları kontrol edin
go version # Go 1.22+

# Linux için derleme
go build -ldflags="-s -w" -o executables/OkulEvrak-Linux .

# Windows için çapraz derleme (Linux üzerinden)
CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc go build -ldflags="-s -w -H=windowsgui" -o executables/OkulEvrak.exe .
```
