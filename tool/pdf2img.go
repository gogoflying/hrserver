
package main

//./main in.pdf 3 4 a.pdf
//unidoc/pdf/model/writer.go +371  s := "Unlicensed UniDoc - Get a license on https://unidoc.io"
//https://github.com/unidoc/unidoc

//2. go get gopkg.in/gographics/imagick.v2/imagick
// sudo apt install libmagic-dev libmagickwand-dev
import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"math"
	"os"

	"github.com/nfnt/resize"
	pdf "github.com/unidoc/unidoc/pdf/model"
	"gopkg.in/gographics/imagick.v2/imagick"
)

func init() {
	// When debugging: use debug-level console logger.
	//unicommon.SetLogger(unicommon.NewConsoleLogger(unicommon.LogLevelDebug))
}

func main() {
	/*sourceFile := "234.png"
	targetFile := "outpdf_6.jpg"
	MergeImage3(targetFile, sourceFile)
	return

	pdfName := "outpdf_6.pdf"
	imageName := "outpdf_6.jpg"

	if err := ConvertPdfToJpg(pdfName, imageName); err != nil {
		fmt.Println(err)
	}
	return*/
	if len(os.Args) < 3 {
		fmt.Printf("Usage: go run pdf_split.go input.pdf <page_from> <page_to> output.pdf\n")
		os.Exit(1)
	}

	inputPath := os.Args[1]
	outputPath := os.Args[2]

	err := splitPdf(inputPath, outputPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Complete, see output file: %s\n", outputPath)
}

func splitPdf(inputPath string, outputPath string) error {

	f, err := os.Open(inputPath)
	if err != nil {
		fmt.Printf("source file not exist err:%v\n", err)
		return err
	}
	defer f.Close()

	pdfReader, err := pdf.NewPdfReader(f)
	if err != nil {
		fmt.Printf("pdf split reader err:%v\n", err)
		return err
	}

	isEncrypted, err := pdfReader.IsEncrypted()
	if err != nil {
		fmt.Printf("pdf check is encrypted,please check err:%v\n", err)
		return err
	}

	if isEncrypted {
		_, err = pdfReader.Decrypt([]byte(""))
		if err != nil {
			return err
		}
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		fmt.Printf("pdf get pages err:%v\n", err)
		return err
	}

	filnameList := []string{}
	for i := 0; i < numPages; i++ {
		name := fmt.Sprintf("outpdf_%d.pdf", i)
		filnameList = append(filnameList, name)
	}

	for item, name := range filnameList {
		pdfWriter := pdf.NewPdfWriter()
		page, err := pdfReader.GetPage(item + 1)
		if err != nil {
			return err
		}
		err = pdfWriter.AddPage(page)
		if err != nil {
			return err
		}

		fWrite, err := os.Create(name)
		if err != nil {
			return err
		}
		defer fWrite.Close()

		err = pdfWriter.Write(fWrite)
		if err != nil {
			return err
		}

		//pdf Convert jpg
		picName := name[0 : len(name)-3]
		picName += "jpg"
		err = ConvertPdfToJpg(name, picName)
		if err != nil {
			fmt.Printf("convertPdfToJpg err:%v\n", err)
			return err
		}
	}

	return nil
}

func ConvertPdfToJpg(pdfName string, imageName string) error {
	// Setup
	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	// Must be *before* ReadImageFile
	// Make sure our image is high quality
	if err := mw.SetResolution(300, 300); err != nil {
		return err
	}

	// Load the image file into imagick
	if err := mw.ReadImage(pdfName); err != nil {
		return err
	}

	// Must be *after* ReadImageFile
	// Flatten image and remove alpha channel, to prevent alpha turning black in jpg
	if err := mw.SetImageAlphaChannel(imagick.ALPHA_CHANNEL_FLATTEN); err != nil {
		return err
	}

	// Set any compression (100 = max quality)
	if err := mw.SetCompressionQuality(95); err != nil {
		return err
	}

	// Select only first page of pdf
	mw.SetIteratorIndex(0)

	// Convert into JPG
	if err := mw.SetFormat("jpg"); err != nil {
		return err
	}

	// Save File
	return mw.WriteImage(imageName)
}

const MaxWidth float64 = 600

func fixSize(img1W, img2W int) (new1W, new2W int) {
	var ( //为了方便计算，将两个图片的宽转为 float64
		img1Width, img2Width = float64(img1W), float64(img2W)
		ratio1, ratio2       float64
	)

	minWidth := math.Min(img1Width, img2Width) // 取出两张图片中宽度最小的为基准

	if minWidth > 600 { // 如果最小宽度大于600，那么两张图片都需要进行缩放
		ratio1 = MaxWidth / img1Width // 图片1的缩放比例
		ratio2 = MaxWidth / img2Width // 图片2的缩放比例

		// 原宽度 * 比例 = 新宽度
		return int(img1Width * ratio1), int(img2Width * ratio2)
	}

	// 如果最小宽度小于600，那么需要将较大的图片缩放，使得两张图片的宽度一致
	if minWidth == img1Width {
		ratio2 = minWidth / img2Width // 图片2的缩放比例
		return img1W, int(img2Width * ratio2)
	}

	ratio1 = minWidth / img1Width // 图片1的缩放比例
	return int(img1Width * ratio1), img2W
}

func MergeImage(soruceImage, targeImage string) {
	file1, _ := os.Open(soruceImage) //打开图片1
	file2, _ := os.Open(targeImage)  //打开图片2
	defer file1.Close()
	defer file2.Close()

	// image.Decode 图片
	var (
		img1, img2 image.Image
		err        error
	)
	if img1, _, err = image.Decode(file1); err != nil {
		log.Fatal(err)
		return
	}
	if img2, _, err = image.Decode(file2); err != nil {
		log.Fatal(err)
		return
	}
	b1 := img1.Bounds()
	b2 := img2.Bounds()
	new1W, new2W := fixSize(b1.Max.X, b2.Max.X)

	// 调用resize库进行图片缩放(高度填0，resize.Resize函数中会自动计算缩放图片的宽高比)
	m1 := resize.Resize(uint(new1W), 0, img1, resize.Lanczos3)
	m2 := resize.Resize(uint(new2W), 0, img2, resize.Lanczos3)

	// 将两个图片合成一张
	newWidth := m1.Bounds().Max.X                                                                          //新宽度 = 随意一张图片的宽度
	newHeight := m1.Bounds().Max.Y + m2.Bounds().Max.Y                                                     // 新图片的高度为两张图片高度的和
	newImg := image.NewNRGBA(image.Rect(0, 0, newWidth, newHeight))                                        //创建一个新RGBA图像
	draw.Draw(newImg, newImg.Bounds(), m1, m1.Bounds().Min, draw.Over)                                     //画上第一张缩放后的图片
	draw.Draw(newImg, newImg.Bounds(), m2, m2.Bounds().Min.Sub(image.Pt(0, m1.Bounds().Max.Y)), draw.Over) //画上第二张缩放后的图片（这里需要注意Y值的起始位置）

	// 保存文件
	imgfile, _ := os.Create("003.jpg")
	defer imgfile.Close()
	jpeg.Encode(imgfile, newImg, &jpeg.Options{100})
}

func MergeImageEx(soruceImage, targeImage string) {
	file, err := os.Create("dst.jpg")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	file1, err := os.Open(soruceImage)
	if err != nil {
		fmt.Println(err)
	}
	defer file1.Close()
	img, _ := jpeg.Decode(file1)

	file2, err := os.Open(targeImage)
	if err != nil {
		fmt.Println(err)
	}
	defer file2.Close()
	img2, _ := jpeg.Decode(file2)

	jpg := image.NewRGBA(image.Rect(200, 200, 500, 500))

	draw.Draw(jpg, jpg.Bounds(), img2, img2.Bounds().Min, draw.Over)                   //首先将一个图片信息存入jpg
	draw.Draw(jpg, jpg.Bounds(), img, img.Bounds().Min.Sub(image.Pt(0, 0)), draw.Over) //将另外一张图片信息存入jpg

	// draw.DrawMask(jpg, jpg.Bounds(), img, img.Bounds().Min, img2, img2.Bounds().Min, draw.Src) // 利用这种方法不能够将两个图片直接合成？目前尚不知道原因。

	jpeg.Encode(file, jpg, nil)

}

func MergeImage3(soruceImage, targeImage string) {
	img_file, err := os.Open(soruceImage)
	if err != nil {
		fmt.Println("打开图片出错")
		fmt.Println(err)
		os.Exit(-1)
	}
	defer img_file.Close()
	img, err := jpeg.Decode(img_file)
	if err != nil {
		fmt.Println("把图片解码为结构体时出错")
		fmt.Println(img)
		os.Exit(-1)
	}

	//水印,用的是我自己支付宝的二维码
	wmb_file, err := os.Open(targeImage)
	if err != nil {
		fmt.Println("打开水印图片出错")
		fmt.Println(err)
		os.Exit(-1)
	}
	defer wmb_file.Close()
	wmb_img, err := png.Decode(wmb_file)
	if err != nil {
		fmt.Println("把水印图片解码为结构体时出错")
		fmt.Println(err)
		os.Exit(-1)
	}

	//把水印写在右下角，并向0坐标偏移10个像素
	offset := image.Pt(img.Bounds().Dx()-wmb_img.Bounds().Dx()-500, img.Bounds().Dy()-wmb_img.Bounds().Dy()-500)
	b := img.Bounds()
	//根据b画布的大小新建一个新图像
	m := image.NewRGBA(b)

	//image.ZP代表Point结构体，目标的源点，即(0,0)
	//draw.Src源图像透过遮罩后，替换掉目标图像
	//draw.Over源图像透过遮罩后，覆盖在目标图像上（类似图层）
	draw.Draw(m, b, img, image.ZP, draw.Src)
	draw.Draw(m, wmb_img.Bounds().Add(offset), wmb_img, image.ZP, draw.Over)

	//生成新图片new.jpg,并设置图片质量
	imgw, err := os.Create("new.jpg")
	jpeg.Encode(imgw, m, &jpeg.Options{100})
	defer imgw.Close()

	fmt.Println("jpg merged success")
}


