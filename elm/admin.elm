module Admin exposing (..)

import Html.App as App
import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (..)
import Task
import Date
import Http
import String
import Dict
import Json.Decode as Json exposing (Decoder, decodeValue, succeed, string, oneOf, null, list, bool, (:=), andThen)
import Json.Encode


main =
    App.program
        { init = init
        , view = view
        , update = update
        , subscriptions = subscriptions
        }


type alias Event =
    { id : Int
    , name : String
    , location : String
    , date : Date.Date
    }


type alias User =
    { name : String
    , email : String
    }


type alias FormFields =
    { email : String
    , password : String
    }


type alias EventDict =
    Dict.Dict Int Event


emptyEvent : Event
emptyEvent =
    { id = 0, name = "", location = "", date = (Date.fromTime 0) }


emptyUser : User
emptyUser =
    { name = "", email = "" }


emptyFormFields : FormFields
emptyFormFields =
    { email = "", password = "" }


type alias Model =
    { loaded : Bool
    , loggedIn : Bool
    , eventInfo : EventDict
    , eventList : List Int
    , user : User
    , form : FormFields
    , curEvt : Event
    , haveEvt : Bool
    }


model : Model
model =
    { loaded = False
    , loggedIn = False
    , eventInfo = (Dict.empty)
    , eventList = []
    , user = emptyUser
    , form = emptyFormFields
    , curEvt = emptyEvent
    , haveEvt = False
    }


init : ( Model, Cmd Msg )
init =
    ( model, loadUser "" )


apiEndpoint =
    "api"


type Msg
    = Noop
    | SetUser (Response User)
    | FailUser Http.Error
    | FormLoginEmail String
    | FormLoginPassword String
    | LogOut
    | DoLogin
    | SetEvents (Response (List Event))
    | SetEvent Event
    | ClearEvent
    | InviteEvent Int
    | DidInviteEvent (Response Bool)
    | FailInviteEvent Http.Error


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        Noop ->
            ( model, Cmd.none )

        DidInviteEvent _ ->
            ( model, Cmd.none )

        FailInviteEvent _ ->
            ( model, Cmd.none )

        SetUser newUser ->
            ( { model | loaded = True, loggedIn = True, user = newUser.data }, loadEvents "" )

        SetEvents evtRes ->
            let
                newEvts =
                    List.map .id evtRes.data
            in
                ( { model | eventList = newEvts, eventInfo = Dict.fromList (List.map makeTupleFromBase evtRes.data) }, Cmd.none )

        FailUser _ ->
            ( { model | loaded = True, loggedIn = False }, Cmd.none )

        FormLoginEmail key ->
            let
                curForm =
                    model.form
            in
                ( { model | form = { curForm | email = key } }, Cmd.none )

        FormLoginPassword key ->
            let
                curForm =
                    model.form
            in
                ( { model | form = { curForm | password = key } }, Cmd.none )

        DoLogin ->
            let
                useDets =
                    model.form
            in
                ( { model | loaded = False }, loginUser useDets )

        SetEvent evt ->
            ( { model | haveEvt = True, curEvt = evt }, Cmd.none)

        ClearEvent ->
            ( { model | haveEvt = False, curEvt = emptyEvent }, Cmd.none)

        InviteEvent evtId ->
            ( model, inviteEvent evtId)

        LogOut ->
            ( model, logOutUser "" )


type alias Response data =
    { errorCode : Int
    , errorMessage : String
    , data : data
    }


dateFromString : String -> Date.Date
dateFromString strDate =
    Date.fromString strDate |> Result.withDefault (Date.fromTime 0)


responseDecoder : Json.Decoder data -> Json.Decoder (Response data)
responseDecoder dataDecoder =
    Json.object3 Response
        ("errorCode" := Json.int)
        ("errorMessage" := Json.string)
        ("data" := dataDecoder)


dateDecoder : Json.Decoder Date.Date
dateDecoder =
    Json.customDecoder string Date.fromString


eventDecoder : Json.Decoder Event
eventDecoder =
    Json.object4 Event
        ("id" := Json.int)
        ("name" := string)
        ("location" := string)
        ("date" := dateDecoder)


userDecoder : Json.Decoder User
userDecoder =
    Json.object2 User
        ("name" := string)
        ("email" := string)


eventsResponseDecoder : Json.Decoder (Response (List Event))
eventsResponseDecoder =
    Json.object3 Response
        ("errorCode" := Json.int)
        ("errorMessage" := string)
        ("data" := (Json.list eventDecoder))


userResponseDecoder : Json.Decoder (Response User)
userResponseDecoder =
    Json.object3 Response
        ("errorCode" := Json.int)
        ("errorMessage" := string)
        ("data" := userDecoder)

boolResponseDecoder : Json.Decoder (Response Bool)
boolResponseDecoder =
    Json.object3 Response
        ("errorCode" := Json.int)
        ("errorMessage" := string)
        ("data" := Json.bool)


makeTupleFromBase : Event -> ( Int, Event )
makeTupleFromBase bse =
    ( bse.id, bse )


loadEvents : String -> Cmd Msg
loadEvents msg =
    let
        url =
            apiEndpoint ++ "/admin/events"
    in
        Task.perform FailUser SetEvents (Http.get eventsResponseDecoder url)


loginUser : FormFields -> Cmd Msg
loginUser form =
    let
        url =
            apiEndpoint ++ "/adminlogin"

        reqJSON =
            Json.Encode.encode 0 (Json.Encode.object [ ( "email", Json.Encode.string form.email ), ( "password", Json.Encode.string form.password ) ])
    in
        Task.perform FailUser SetUser (Http.post userResponseDecoder url (Http.string reqJSON))

inviteEvent : Int -> Cmd Msg
inviteEvent evtId =
    let
        url =
            apiEndpoint ++ "/admin/event/" ++ (toString evtId) ++ "/invite"
    in
        Task.perform FailInviteEvent DidInviteEvent (Http.post boolResponseDecoder url (Http.string ""))

logOutUser : String -> Cmd Msg
logOutUser str =
    let
        url =
            apiEndpoint ++ "/logout"
    in
        Task.perform FailUser SetUser (Http.get userResponseDecoder url)


loadUser : String -> Cmd Msg
loadUser msg =
    let
        url =
            apiEndpoint ++ "/admin/me"
    in
        Task.perform FailUser SetUser (Http.get userResponseDecoder url)


formatDate : Date.Date -> String
formatDate date =
    (toString (Date.day date)) ++ " " ++ (toString (Date.month date)) ++ " " ++ (toString (Date.year date))


drawNotLoggedIn : Model -> Html Msg
drawNotLoggedIn model =
    div [ class "login" ]
        [ div [ class "title" ] [ text "Please log in" ]
        , div [ class "body" ]
            [ div []
                [ span [ class "label" ] [ text "Email" ]
                , input [ type' "text", onInput FormLoginEmail ] []
                ]
            , div []
                [ span [ class "label" ] [ text "Password" ]
                , input [ type' "password", onInput FormLoginPassword ] []
                ]
            , div []
                [ button [ onClick DoLogin ] [ text "Login" ]
                ]
            ]
        ]


drawLoading : Model -> Html a
drawLoading model =
    div []
        [ div [] [ text "Loading" ]
        ]


drawOneEvent : EventDict -> Int -> Html Msg
drawOneEvent evtInfo evtId =
    let
        event =
            Maybe.withDefault emptyEvent (Dict.get evtId evtInfo)
    in
        drawEvent event

drawEvent : Event -> Html Msg
drawEvent event =
        div [ class "event" ]
            [ div [ class "name" ] [ a [ onClick (SetEvent event) ] [text event.name ] ]
            , div [ class "location" ] [ text event.location ]
            , div [ class "date" ] [ text (formatDate event.date) ]
            , div [ class "actions" ]
                [ a [ class "getList", href (apiEndpoint ++ "/admin/event/" ++ (toString event.id) ++ "/summary"), target "new" ] [ text "Get Submissions" ]
                ]
            ]


view : Model -> Html Msg
view model =
    if model.loaded then
        if model.loggedIn then
            if model.haveEvt then
                    div [ class "loaded" ]
                        [ div [ id "header" ]
                            [ div [ class "user" ] [ span [ class "name" ] [ text model.user.name ] ]
                            ]
                        , div [ id "body" ]
                            [ div [ id "eventIntro" ] [ (text "These are the events you're organising:") ]
                            , div [ id "events" ] [ drawEvent model.curEvt ]
                            , div [ id "actions" ] [
                                    a [ onClick ( InviteEvent model.curEvt.id ) ] [ text "Send invitations" ]
                                ]
                            ]
                        , div [ id "footer" ]
                            [ a [ onClick ClearEvent ] [ text "Back to list" ]
                            ]
                        ]

            else
                let
                    evtRendr =
                        drawOneEvent model.eventInfo
                in
                    div [ class "loaded" ]
                        [ div [ id "header" ]
                            [ div [ class "user" ] [ span [ class "name" ] [ text model.user.name ] ]
                            ]
                        , div [ id "body" ]
                            [ div [ id "eventIntro" ] [ (text "These are the events you're organising:") ]
                            , div [ id "events" ] (List.map evtRendr model.eventList)
                            ]
                        , div [ id "footer" ]
                            [ a [ onClick LogOut ] [ text "Log out" ]
                            ]
                        ]
        else
            drawNotLoggedIn model
    else
        drawLoading model


subscriptions : Model -> Sub Msg
subscriptions model =
    Sub.none
