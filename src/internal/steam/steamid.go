/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package steam

import (
  "strings"
  "strconv"
  "errors"
)

type SteamID struct {
  Universe  uint64
  Type      uint64
  Instance  uint64
  AccountID uint64
}

const (
  EUniverseInvalid  = 0
  EUniversePublic   = 1
  EUniverseBeta     = 2
  EUniverseInternal = 3
  EUniverseDev      = 4
)

const (
  ETypeInvalid      = 0
  ETypeIndividual   = 1
  ETypeMultiseat    = 2
  ETypeGameServer   = 3
  ETypeAnonGameSrv  = 4
  ETypePending      = 5
  ETypeContentSrv   = 6
  ETypeClan         = 7
  ETypeChat         = 8
  ETypeP2P          = 9
)

const (
  EInstanceAll      = 0
  EInstanceDesktop  = 1
  EInstanceConsole  = 2
  EInstanceWeb      = 4
)

func FromSteam64(id uint64) *SteamID {
  return &SteamID{
    Universe:  (id >> 56) & 0xFF,
    Type:      (id >> 52) & 0xF,
    Instance:  (id >> 32) & 0xFFFFF,
    AccountID: id & 0xFFFFFFFF,
  }
}

// "STEAM_X:Y:Z"
func FromSteam2(s string) (*SteamID, error) {
  parts := strings.Split(s, ":")
  if len(parts) != 3 || !strings.HasPrefix(parts[0], "STEAM_") {
    return nil, errors.New("Invalid Steam2 ID")
  }
  X, _ := strconv.ParseUint(parts[0][6:], 10, 64) 
  Y, _ := strconv.ParseUint(parts[1], 10, 64)
  Z, _ := strconv.ParseUint(parts[2], 10, 64)

  var universe uint64 = EUniversePublic
  if X > 1 {
    universe = X
  }
  
  return &SteamID{
    Universe: universe, 
    Type: ETypeIndividual, 
    Instance: EInstanceDesktop, 
    AccountID: Z*2 + Y,
  }, nil
}

// "[U:1:Z]"
func FromSteam3(s string) (*SteamID, error) {
  s = strings.Trim(s, "[]")
  parts := strings.Split(s, ":")
  if len(parts) != 3 {
    return nil, errors.New("Invalid Steam3 ID")
  }

  universe, err := strconv.ParseUint(parts[1], 10, 64);
    if err != nil {
    return nil, err
  }
  
  accountID, err := strconv.ParseUint(parts[2], 10, 64)
  if err != nil {
    return nil, err
  }
  
  sid := &SteamID{
    Universe: universe, 
    Type: ETypeIndividual, 
    Instance: EInstanceAll, 
    AccountID: accountID,
  }

  switch parts[0] {
    case "I": 
      sid.Type = ETypeInvalid
    case "U":
      sid.Type = ETypeIndividual
      sid.Instance = EInstanceDesktop
    case "M":
      sid.Type = ETypeMultiseat
    case "G":
      sid.Type = ETypeGameServer
    case "A":
      sid.Type = ETypeAnonGameSrv
    case "P":
      sid.Type = ETypePending
    case "C":
      sid.Type = ETypeContentSrv
    case "g":
      sid.Type = ETypeClan
    case "c":
      sid.Type = ETypeChat
    case "a":
      sid.Type = ETypeP2P
    default:
      sid.Type = ETypeIndividual
  }

  return sid, nil
}

func ParseSteamID(input string) (*SteamID, error) {
  if strings.HasPrefix(input, "STEAM_") {
    return FromSteam2(input)
  }
  if strings.HasPrefix(input, "[") && strings.HasSuffix(input, "]") {
    return FromSteam3(input)
  }
  if id, err := strconv.ParseUint(input, 10, 64); err == nil {
    return FromSteam64(id), nil
  }
  return nil, errors.New("Unknown SteamID format")
}

func (sid *SteamID) AsSteam64() string {
  id := ((sid.Universe & 0xFF) << 56) |
        ((sid.Type & 0xF) << 52) |
        ((sid.Instance & 0xFFFFF) << 32) |
        (sid.AccountID & 0xFFFFFFFF)
  return strconv.FormatUint(id, 10)
}

func (sid *SteamID) AsSteam3() string {
  result := "["
  switch sid.Type {
    case ETypeInvalid:
      result += "I"
    case ETypeIndividual:
      result += "U"
    case ETypeMultiseat:
      result += "M"
    case ETypeGameServer:
      result += "G"
    case ETypeAnonGameSrv:
      result += "A"
    case ETypePending:
      result += "P"
    case ETypeContentSrv:
      result += "C"
    case ETypeClan:
      result += "g"
    case ETypeChat:
      result += "c"
    case ETypeP2P:
      result += "a"
    default:
      result += "U"
  }
  result += ":"
  result += strconv.FormatUint(sid.Universe, 10)
  result += ":"
  result += strconv.FormatUint(sid.AccountID, 10)
  result += "]"

  return result
}

func (sid *SteamID) AsSteam2() string {
  Y := sid.AccountID % 2
  Z := sid.AccountID / 2

  result := "STEAM_"
  result += strconv.FormatUint(sid.Universe, 10)
  result += ":"
  result += strconv.FormatUint(Y, 10)
  result += ":"
  result += strconv.FormatUint(Z, 10)
  
  return result 
}