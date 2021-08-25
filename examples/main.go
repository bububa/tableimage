package main

import (
	"image"
	"log"
	"os"

	"github.com/bububa/tableimage"
	"github.com/llgcode/draw2d"
)

func main() {
	imageURL := "https://images.pexels.com/photos/906052/pexels-photo-906052.jpeg?auto=compress&cs=tinysrgb&dpr=2&h=750&w=200"

	ti, err := tableimage.New(
		tableimage.WithBgColor("#FFFFFF"),
		tableimage.WithBorderColor("#0277BD"),
		tableimage.WithFontSize(11),
		tableimage.WithDPI(144),
		tableimage.WithFontData(&draw2d.FontData{
			Name:   "NotoSansCJKsc",
			Family: draw2d.FontFamilySans,
			Style:  draw2d.FontStyleNormal,
		}),
		tableimage.WithFontFolder("./font"),
	)
	if err != nil {
		log.Fatalln(err)
		return
	}

	headerFont := &tableimage.Font{
		Size: 13,
		Data: &draw2d.FontData{
			Name:   "Roboto",
			Family: draw2d.FontFamilySans,
			Style:  draw2d.FontStyleBold,
		},
	}

	footerFont := &tableimage.Font{
		Size: 10,
		Data: &draw2d.FontData{
			Name:   "NotoSansCJKsc",
			Family: draw2d.FontFamilySans,
			Style:  draw2d.FontStyleNormal,
		},
	}

	headerStyle := &tableimage.Style{
		Font: headerFont,
	}

	caption := &tableimage.Cell{
		Text: "this is a caption very long long long ago adsfasdfwe;alsdkfasdfasf asdfajsdf",
		Style: &tableimage.Style{
			Color: "#3F51B5",
			Font:  headerFont,
		},
	}

	tfooter := &tableimage.Cell{
		Text: "this is a tabel footer very long long long ago adsfasdfwe;alsdkfasdfasf asdfajsdf",
		Style: &tableimage.Style{
			Color: "#D7CCC8",
			Font:  footerFont,
		},
	}

	rows := []tableimage.Row{
		{
			Style: headerStyle,
			Cells: []tableimage.Cell{
				{
					Text: "Id",
				},
				{
					Text: "Name",
				},
				{
					Style: &tableimage.Style{
						Color: "#008000",
					},
					Text: "Price",
				},
			},
		},
		{
			Cells: []tableimage.Cell{
				{
					Text: "2223",
				},
				{
					Style: &tableimage.Style{
						Color:    "#000",
						MaxWidth: 100,
						Align:    tableimage.CENTER,
					},
					Text: "这是一个真正的table图片",
				},
				{
					Style: &tableimage.Style{
						Color: "#0000ff",
					},
					Text: "2000$",
				},
			},
		},
		{
			Cells: []tableimage.Cell{
				{
					Style: &tableimage.Style{
						Color:  "#6A1B9A",
						VAlign: tableimage.TOP,
					},
					Text: "11",
				},
				{
					Style: &tableimage.Style{
						Color:    "#FFF",
						MaxWidth: 100,
						BgColor:  "#D32F2F",
					},
					IgnoreInlineStyle: false,
					Text:              "A more <text bgcolor='#8BC34A' color='#000' padding='4'>cooler product this</text> time on 3 lines",
				},
				{
					Style: &tableimage.Style{
						Color:  "#0000ff",
						Align:  tableimage.RIGHT,
						VAlign: tableimage.BOTTOM,
					},
					Text: "200$",
				},
			},
		},
		{
			Cells: []tableimage.Cell{
				{
					Text: "2223",
					Image: &tableimage.Image{
						URL:     imageURL,
						Size:    image.Pt(80, 0),
						VAlign:  tableimage.BOTTOM,
						Padding: tableimage.NewPaddingY(4),
					},
					Style: &tableimage.Style{
						Align: tableimage.RIGHT,
					},
				},
				{
					Style: &tableimage.Style{
						Color:    "#000",
						MaxWidth: 100,
						Align:    tableimage.CENTER,
					},
					Image: &tableimage.Image{
						URL:     imageURL,
						Size:    image.Pt(80, 0),
						VAlign:  tableimage.TOP,
						Padding: tableimage.NewPaddingY(4),
					},
					Text: "这是一个真正的table图片",
				},
				{
					Style: &tableimage.Style{
						Color: "#0000ff",
					},
					Text: "2000$",
				},
			},
		},
		{
			Cells: []tableimage.Cell{
				{
					Text: "2223",
					Image: &tableimage.Image{
						URL:     imageURL,
						Size:    image.Pt(80, 0),
						Align:   tableimage.LEFT,
						Padding: tableimage.NewPaddingX(4),
					},
					Style: &tableimage.Style{
						VAlign: tableimage.BOTTOM,
					},
				},
				{
					Style: &tableimage.Style{
						Color:    "#000",
						MaxWidth: 150,
						Align:    tableimage.RIGHT,
						VAlign:   tableimage.MIDDLE,
					},
					Image: &tableimage.Image{
						URL:     imageURL,
						Size:    image.Pt(80, 0),
						Align:   tableimage.RIGHT,
						Padding: tableimage.NewPaddingX(4),
					},
					Text: "这是一个真正的table图片",
				},
				{
					Style: &tableimage.Style{
						Color: "#0000ff",
					},
					Text: "2000$",
				},
			},
		},
	}
	img, err := ti.Draw(rows, caption, tfooter)
	if err != nil {
		log.Fatalln(err)
		return
	}
	f, err := os.Create("./test.png")
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer f.Close()
	tableimage.Write(f, img, tableimage.PNG)
}
