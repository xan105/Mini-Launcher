/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package ui

import (
  "golang.org/x/sys/windows"
)

func getScreenResolution() (uint32, uint32, error){
  hDC, err := getDC(0)
  if err != nil {
    return 0, 0, err
  }
  defer releaseDC(0, hDC)
  
  width := getDeviceCaps(hDC, HORZRES)
  height := getDeviceCaps(hDC, VERTRES)
  return width, height, nil
}

func createBrushFromBMP(splashImage string) (windows.Handle, BITMAP, error) {
  hbm, err := loadImage(splashImage)
  if err != nil {    
    return 0, BITMAP{}, err
  } 
   
  hbrush, err := createPatternBrush(hbm)
  if err != nil {    
    return 0, BITMAP{}, err
  }
  
  //Get Image dimension
  image, err:= getObject(hbm)
  if err != nil {
    return hbrush, image, err
  }

  return hbrush, image, nil
}