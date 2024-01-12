package fmtx

import (
	"fmt"
	"io"
)

func Lorem(w io.Writer) {
	fmt.Fprint(w, "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Fusce mattis nulla eget rutrum viverra. Fusce hendrerit neque consectetur neque tempor, nec ornare mi vestibulum. Donec gravida condimentum dolor, at fermentum lectus molestie ultrices. Aliquam tincidunt diam in magna venenatis rhoncus. Proin pellentesque ante id neque varius, in aliquam arcu semper. Nunc eget risus eget elit venenatis vehicula. Proin bibendum, nisl et sodales pharetra, sem tellus placerat odio, nec posuere est neque vel tortor. Nunc vel dignissim orci. Donec risus ex, condimentum ac congue vitae, vestibulum id risus. Cras rutrum vulputate sodales. Cras ultrices nisi vitae velit sodales elementum. Sed eu enim nunc. Vivamus venenatis, ligula vitae commodo convallis, felis orci posuere tortor, ut posuere ex lectus a mi. Vivamus dictum scelerisque risus, non vestibulum enim dictum sed. Nulla convallis, leo volutpat consectetur ultricies, tellus ex suscipit est, vitae suscipit arcu erat eu sem.\n")
	fmt.Fprint(w, "\n")
	fmt.Fprint(w, "Nulla eget ligula nec risus varius vulputate vitae non odio. Nullam pretium congue justo at vestibulum. Morbi a ipsum turpis. Fusce in elit nibh. Morbi ut pulvinar est. Ut scelerisque hendrerit nunc in faucibus. Vivamus ante neque, finibus in arcu ac, sagittis placerat libero. Integer ut felis a est finibus euismod nec id ex. Donec vel felis nisl. Duis a vehicula lectus, quis tincidunt turpis. Pellentesque in libero sit amet magna vestibulum vestibulum. Duis interdum egestas suscipit. Nulla lectus ligula, fringilla pellentesque pulvinar quis, aliquet vitae lectus. Proin metus ligula, sagittis vel quam id, viverra posuere ante. Duis accumsan sit amet est tincidunt vehicula. Nam ante dui, semper non tortor quis, pulvinar tincidunt justo.\n")
}
