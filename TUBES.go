package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)


 

// Tipe Bentukan
type Waktu struct {
	JamMasuk  time.Time
	JamKeluar time.Time
}

type Kendaraan struct {
	PlatNomor string
	Jenis     string
	Waktu     Waktu
	Slot      int
}

type SlotParkir struct {
	Nomor  int
	Kosong bool
}

// Variabel global
var slotParkir [100]SlotParkir
var kendaraanParkir []Kendaraan
var historiKendaraan []Kendaraan
var scanner = bufio.NewScanner(os.Stdin)

// Inisialisasi slot
func initSlot() {
	for i := 0; i < len(slotParkir); i++ {
		slotParkir[i] = SlotParkir{Nomor: i + 1, Kosong: true}
	}
}

// Fungsi input
func input(prompt string) string {
	fmt.Print(prompt)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

// Fungsi: Hitung durasi parkir
func hitungDurasi(waktu Waktu) time.Duration {
	return waktu.JamKeluar.Sub(waktu.JamMasuk)
}

// Prosedur: Masukkan kendaraan dengan validasi panjang plat nomor dan input nomor slot
func masukkanKendaraan() {
	const maxPlatLength = 10

	plat := input("Masukkan plat nomor (maks 10 karakter): ")
	if len(plat) > maxPlatLength {
		fmt.Printf("❌ Plat nomor terlalu panjang! Maksimum %d karakter.\n", maxPlatLength)
		return
	}

	jenis := input("Masukkan jenis kendaraan (Mobil/Motor): ")
	slotInput := input("Masukkan nomor slot yang diinginkan: ")
	slotNum, err := strconv.Atoi(slotInput)
	if err != nil || slotNum < 1 || slotNum > len(slotParkir) {
		fmt.Println("Nomor slot tidak valid.")
		return
	}
	if !slotParkir[slotNum-1].Kosong {
		fmt.Println("Slot sudah terisi.")
		return
	}

	now := time.Now()
	slotParkir[slotNum-1].Kosong = false
	kendaraan := Kendaraan{
		PlatNomor: plat,
		Jenis:     jenis,
		Waktu:     Waktu{JamMasuk: now},
		Slot:      slotNum,
	}
	kendaraanParkir = append(kendaraanParkir, kendaraan)
	fmt.Println("✅ Kendaraan masuk ke slot:", slotNum)
}

// Prosedur: Keluarkan kendaraan dan simpan histori
func keluarkanKendaraan() {
	plat := input("Masukkan plat nomor kendaraan yang keluar: ")
	for i, k := range kendaraanParkir {
		if k.PlatNomor == plat {
			now := time.Now()
			k.Waktu.JamKeluar = now
			durasi := hitungDurasi(k.Waktu)

			historiKendaraan = append(historiKendaraan, k)
			slotParkir[k.Slot-1].Kosong = true
			kendaraanParkir = append(kendaraanParkir[:i], kendaraanParkir[i+1:]...)

			fmt.Printf("Kendaraan keluar dari slot: %d\n", k.Slot)
			fmt.Printf("Jenis: %s\n", k.Jenis)
			fmt.Printf("Durasi parkir: %.0f menit\n", durasi.Minutes())

			return
		}
	}
	fmt.Println("❌ Kendaraan tidak ditemukan!")
}

// Sequential Search: Cari kendaraan berdasarkan plat
func CariKendaraanSequential() {
	plat := input("Masukkan plat nomor: ")
	ditemukan := false

	for i := 0; i < len(kendaraanParkir); i++ {
		if kendaraanParkir[i].PlatNomor == plat {
			fmt.Printf("Ditemukan: %s (%s), Slot %d\n", kendaraanParkir[i].PlatNomor, kendaraanParkir[i].Jenis, kendaraanParkir[i].Slot)
			ditemukan = true
		}
	}

	if !ditemukan {
		fmt.Println("Kendaraan tidak ditemukan.")
	}
}

// Binary Search: Cari kendaraan berdasarkan jam masuk (HH:MM), tanpa menggunakan break
func cariKendaraanBerdasarkanJam() {
	startStr := input("Masukkan jam mulai (HH:MM): ")
	endStr := input("Masukkan jam akhir (HH:MM): ")

	// Parse input jam menjadi total menit
	parseJam := func(s string) int {
		parts := strings.Split(s, ":")
		if len(parts) != 2 {
			return -1
		}
		jam, _ := strconv.Atoi(parts[0])
		menit, _ := strconv.Atoi(parts[1])
		return jam*60 + menit
	}

	startMenit := parseJam(startStr)
	endMenit := parseJam(endStr)

	if startMenit == -1 || endMenit == -1 {
		fmt.Println("Format jam tidak valid. Gunakan format HH:MM.")
		return
	}

	ditemukan := false
	fmt.Printf("Kendaraan yang masuk antara %s dan %s:\n", startStr, endStr)

	for _, k := range kendaraanParkir {
		jamMasuk := k.Waktu.JamMasuk
		jam := jamMasuk.Hour()
		menit := jamMasuk.Minute()
		totalMenit := jam*60 + menit

		if totalMenit >= startMenit && totalMenit <= endMenit {
			fmt.Printf("- %s (%s), Slot %d, Masuk: %s\n",
				k.PlatNomor, k.Jenis, k.Slot, jamMasuk.Format("15:04"))
			ditemukan = true
		}
	}

	if !ditemukan {
		fmt.Println("Tidak ada kendaraan dalam rentang waktu tersebut.")
	}
}

// Opsi 4: Tampilkan daftar slot kosong tanpa input apapun
func cariSlotKosong() {
	fmt.Println("Daftar slot kosong:")
	var kosong []SlotParkir
	for _, slot := range slotParkir {
		if slot.Kosong {
			kosong = append(kosong, slot)
		}
	}
	if len(kosong) == 0 {
		fmt.Println("Tidak ada slot kosong.")
		return
	}
	for _, s := range kosong {
		fmt.Printf("Slot %d\n", s.Nomor)
	}
}

// Selection Sort: Urutkan histori berdasarkan durasi dan tampilkan
func urutkanKendaraanParkirBerdasarkanDurasi() {
	if len(kendaraanParkir) == 0 {

		fmt.Println("Tidak ada kendaraan yang sedang parkir.")

		return
	}

	// Selection Sort langsung di slice asli
	for i := 0; i < len(kendaraanParkir); i++ {
		min := i
		for j := i + 1; j < len(kendaraanParkir); j++ {
			durasiJ := time.Since(kendaraanParkir[j].Waktu.JamMasuk)
			durasiMin := time.Since(kendaraanParkir[min].Waktu.JamMasuk)
			if durasiJ < durasiMin {
				min = j
			}
		}
		kendaraanParkir[i], kendaraanParkir[min] = kendaraanParkir[min], kendaraanParkir[i]
	}

	// Tampilkan hasil

	fmt.Println("Kendaraan parkir diurutkan berdasarkan durasi parkir hingga saat ini:")
	for _, k := range kendaraanParkir {
		durasi := time.Since(k.Waktu.JamMasuk)
		fmt.Printf("- %s (%s), Slot %d, Durasi: %.0f menit\n",
			k.PlatNomor, k.Jenis, k.Slot, durasi.Minutes())
	}
}

// Binary Sort: Urutkan kendaraan berdasarkan waktu masuk dan tampilkan
func urutkanHistoriBerdasarkanJenisDanJamKeluar() {
	if len(historiKendaraan) == 0 {
		fmt.Println("📭 Belum ada histori kendaraan.")
		return
	}

	// Pisahkan histori berdasarkan jenis
	var motorList []Kendaraan
	var mobilList []Kendaraan

	for _, k := range historiKendaraan {
		jenisLower := strings.ToLower(k.Jenis)
		if jenisLower == "motor" {
			motorList = append(motorList, k)
		} else {
			mobilList = append(mobilList, k)
		}
	}

	// Fungsi sorting menggunakan Insertion Sort berdasarkan JamKeluar (ascending)
	insertionSortByJamKeluar := func(list []Kendaraan) {
		for i := 1; i < len(list); i++ {
			key := list[i]
			j := i - 1
			for j >= 0 && list[j].Waktu.JamKeluar.After(key.Waktu.JamKeluar) {
				list[j+1] = list[j]
				j--
			}
			list[j+1] = key
		}
	}

	insertionSortByJamKeluar(motorList)
	insertionSortByJamKeluar(mobilList)

	// Tampilkan hasil

	fmt.Println("Histori kendaraan diurutkan berdasarkan jenis dan waktu keluar:")

	if len(motorList) > 0 {
		fmt.Println("- Motor:")
		for _, m := range motorList {
			fmt.Printf("  - %s, Slot %d, Keluar: %s\n",
				m.PlatNomor, m.Slot, m.Waktu.JamKeluar.Format("15:04:05"))
		}
	} else {
		fmt.Println("- Motor: Tidak ada")
	}

	if len(mobilList) > 0 {
		fmt.Println("- Mobil:")
		for _, m := range mobilList {
			fmt.Printf("  - %s, Slot %d, Keluar: %s\n",
				m.PlatNomor, m.Slot, m.Waktu.JamKeluar.Format("15:04:05"))
		}
	} else {
		fmt.Println("- Mobil: Tidak ada")
	}
}

// Menampilkan kendaraan yang sedang parkir
func tampilkanKendaraanParkir() {
	if len(kendaraanParkir) == 0 {
		fmt.Println("Tidak ada kendaraan yang sedang parkir.")
		return
	}
	fmt.Println("Daftar kendaraan parkir:")
	for _, k := range kendaraanParkir {
		fmt.Printf("- %s (%s), Slot %d, Masuk: %s\n", k.PlatNomor, k.Jenis, k.Slot, k.Waktu.JamMasuk.Format("15:04:05"))
	}
}

// Menampilkan histori kendaraan
func tampilkanHistori() {
	if len(historiKendaraan) == 0 {
		fmt.Println("Belum ada histori kendaraan.")
		return
	}
	fmt.Println("Histori kendaraan:")
	for _, h := range historiKendaraan {
		durasi := hitungDurasi(h.Waktu)
		fmt.Printf("- %s (%s), Slot %d, Durasi: %.0f menit\n",
			h.PlatNomor, h.Jenis, h.Slot, durasi.Minutes())
	}
}

func main() {
	initSlot()
	for {
		fmt.Println("\n===== MENU PARKIR =====")
		fmt.Println("1. Masukkan Kendaraan")
		fmt.Println("2. Keluarkan Kendaraan")
		fmt.Println("3. Cari Kendaraan (Sequential Search)")
		fmt.Println("4. Cari Kendaraan Berdasarkan Waktu (Binary Search)")
		fmt.Println("5. Cari Slot Kosong (Sequential Search)")
		fmt.Println("6. Tampilkan Kendaraan yang Parkir)")
		fmt.Println("7. Tampilkan Histori Kendaraan")
		fmt.Println("8. Urutkan Riwayat Kendaraan Berdasarkan Durasi Parkir (Selection Sort)")
		fmt.Println("9. Urutkan Histori Berdasarkan Jenis dan Waktu Keluar (Insertion Sort)")

		fmt.Println("0. Keluar")
		pilihan := input("Pilih menu: ")

		switch pilihan {
		case "1":
			masukkanKendaraan()
		case "2":
			keluarkanKendaraan()
		case "3":
			CariKendaraanSequential()
		case "4":
			cariKendaraanBerdasarkanJam()
		case "5":
			cariSlotKosong()
		case "6":
			tampilkanKendaraanParkir()
		case "7":
			tampilkanHistori()
		case "8":
			urutkanKendaraanParkirBerdasarkanDurasi()
		case "9":
			urutkanHistoriBerdasarkanJenisDanJamKeluar()
		case "0":
			fmt.Println("Terima kasih telah menggunakan sistem parkir!")
			return
		default:
			fmt.Println("Pilihan tidak valid.")
		}
	}
}