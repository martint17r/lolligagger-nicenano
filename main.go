package main

import (
	"log"

	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

// material shrinkage
const shrink = 1.0 / 0.999 // PLA ~0.1%
//const shrink = 1.0/0.995; // ABS ~0.5%

//-----------------------------------------------------------------------------

func holder() (sdf.SDF3, error) {
	// dimensions taken from nice!nano v1
	width := 18.0
	length := 33.47

	trayHeight := 5.0
	pinClearance := 2.8
	pinOffset := 0.15
	wallThickness := 1.0

	slotHeight := 12.0
	slotWidth := 28.0
	slotLength := 3.9 + 2
	slotOffsetY := -2.0

	shieldHeight := 15.0
	shieldWidth := 31.0
	shieldLength := 1.6

	trayBottomHeight := 1.5

	middlePinWidth := 9.0
	middlePinPos := 26.8

	tray, err := sdf.Box3D(sdf.V3{
		X: width + 2*wallThickness,
		Y: length + 2*wallThickness,
		Z: trayHeight,
	}, 0)

	if err != nil {
		return nil, err
	}

	tray = sdf.Transform3D(tray, sdf.Translate3d(sdf.V3{
		X: 0,
		Y: 0,
		Z: tray.BoundingBox().Size().Z / 2}))

	slot, err := sdf.Box3D(sdf.V3{
		X: slotWidth,
		Y: slotLength,
		Z: slotHeight,
	}, 0)

	if err != nil {
		return nil, err
	}

	slot = sdf.Transform3D(slot, sdf.Translate3d(sdf.V3{
		X: 0,
		Y: slotLength/2 + slotOffsetY + (length+2*wallThickness)/2,
		Z: slotHeight / 2}))

	tray = sdf.Union3D(tray, slot)

	nnBB, err := sdf.Box3D(sdf.V3{
		X: width,
		Y: length,
		Z: trayHeight,
	}, 0)
	nnBB = sdf.Transform3D(nnBB, sdf.Translate3d(sdf.V3{
		X: 0,
		Y: wallThickness,
		Z: nnBB.BoundingBox().Size().Z/2 + trayBottomHeight}))

	if err != nil {
		return nil, err
	}
	tray = sdf.Difference3D(tray, nnBB)

	tray = sdf.Transform3D(tray, sdf.Translate3d(sdf.V3{
		X: 0,
		Y: -(length + 2*wallThickness) / 2,
		Z: 0}))

	pinBox, err := sdf.Box3D(sdf.V3{
		X: pinClearance,
		Y: length,
		Z: trayHeight,
	}, 0)

	if err != nil {
		return nil, err
	}

	pinBoxLeft := sdf.Transform3D(pinBox, sdf.Translate3d(sdf.V3{
		X: -width/2 + pinClearance/2,
		Y: wallThickness - pinBox.BoundingBox().Size().Y/2,
		Z: pinOffset + pinBox.BoundingBox().Size().Z/2}))

	tray = sdf.Difference3D(tray, pinBoxLeft)

	pinBoxRight := sdf.Transform3D(pinBox, sdf.Translate3d(sdf.V3{
		X: width/2 - pinClearance/2,
		Y: wallThickness - pinBox.BoundingBox().Size().Y/2,
		Z: pinOffset + pinBox.BoundingBox().Size().Z/2}))

	tray = sdf.Difference3D(tray, pinBoxRight)

	pinMiddleBox, err := sdf.Box3D(sdf.V3{
		X: middlePinWidth,
		Y: pinClearance,
		Z: trayHeight,
	}, 0)

	if err != nil {
		return nil, err
	}

	pinMiddleBox = sdf.Transform3D(pinMiddleBox, sdf.Translate3d(sdf.V3{
		X: -middlePinWidth / 2,
		Y: -middlePinPos,
		Z: pinOffset + trayHeight/2}))

	tray = sdf.Difference3D(tray, pinMiddleBox)
	/*
		tray = sdf.Transform3D(tray, sdf.Translate3d(sdf.V3{
			X: 5,
			Y: 0,
			Z: 0}))
	*/

	shield, err := sdf.Box3D(sdf.V3{
		X: shieldWidth,
		Y: shieldLength,
		Z: shieldHeight,
	}, 0)

	shield = sdf.Transform3D(shield, sdf.Translate3d(sdf.V3{
		X: 0,
		Y: shieldLength/2 + slotLength + slotOffsetY,
		Z: shieldHeight / 2}))

	tray = sdf.Union3D(tray, shield)

	slotCutoutThickness := 0.5

	usbCutout, err := sdf.Box3D(sdf.V3{
		X: 11.5,
		Y: slotLength + shieldLength - slotCutoutThickness,
		Z: 11 + trayBottomHeight,
	}, 0)

	if err != nil {
		return nil, err
	}

	usbCutout = sdf.Transform3D(usbCutout, sdf.Translate3d(sdf.V3{
		X: 0,
		Y: slotCutoutThickness + usbCutout.BoundingBox().Size().Y/2,
		Z: 0}))
	tray = sdf.Difference3D(tray, usbCutout)

	usbHole, err := sdf.Box3D(sdf.V3{
		X: 9,
		Y: slotLength + shieldLength,
		Z: 3.2,
	}, 0)

	if err != nil {
		return nil, err
	}

	usbHole = sdf.Transform3D(usbHole, sdf.Translate3d(sdf.V3{
		X: 0,
		Y: usbHole.BoundingBox().Size().Y / 2,
		Z: trayBottomHeight + usbHole.BoundingBox().Size().Z/2}))
	tray = sdf.Difference3D(tray, usbHole)

	switchCutout, err := sdf.Box3D(sdf.V3{
		X: 12,
		Y: 5.3,
		Z: 4.54,
	}, 0)

	if err != nil {
		return nil, err
	}

	swCoHeight := 9.2
	switchCutout = sdf.Transform3D(switchCutout, sdf.Translate3d(sdf.V3{
		X: 0,
		Y: slotCutoutThickness + switchCutout.BoundingBox().Size().Y/2,
		Z: swCoHeight}))
	tray = sdf.Difference3D(tray, switchCutout)

	switchHole, err := sdf.Box3D(sdf.V3{
		X: 7.75,
		Y: slotLength + shieldLength,
		Z: 1.2,
	}, 0)

	if err != nil {
		return nil, err
	}

	switchHole = sdf.Transform3D(switchHole, sdf.Translate3d(sdf.V3{
		X: 0,
		Y: switchHole.BoundingBox().Size().Y/2 + slotOffsetY,
		Z: swCoHeight}))
	tray = sdf.Difference3D(tray, switchHole)

	tL := sdf.NewPolygon()
	tL.Add(0, 0)
	tL.Add(0, 2)
	tL.Add(-2, 0)
	tL.Close()

	tpL, err := sdf.Polygon2D(tL.Vertices())
	if err != nil {
		return nil, err
	}
	triangleL := sdf.Extrude3D(tpL, slotHeight)
	triangleL = sdf.Transform3D(triangleL, sdf.Translate3d(sdf.V3{
		X: -slotWidth / 2,
		Y: -triangleL.BoundingBox().Size().Y,
		Z: triangleL.BoundingBox().Size().Z / 2}))
	tray = sdf.Union3D(tray, triangleL)

	tR := sdf.NewPolygon()
	tR.Add(0, 0)
	tR.Add(0, -2)
	tR.Add(2, -2)
	tR.Close()

	tpR, err := sdf.Polygon2D(tR.Vertices())
	if err != nil {
		return nil, err
	}
	triangleR := sdf.Extrude3D(tpR, slotHeight)
	triangleR = sdf.Transform3D(triangleR, sdf.Translate3d(sdf.V3{
		X: +slotWidth / 2,
		Y: 0,
		Z: triangleR.BoundingBox().Size().Z / 2}))
	tray = sdf.Union3D(tray, triangleR)

	return tray, nil
}

func main() {
	s, err := holder()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.RenderSTL(sdf.ScaleUniform3D(s, shrink), 300, "holder.stl")
}

//-----------------------------------------------------------------------------
