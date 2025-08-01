const UNAUTHORIZED_STATUS = 0
const AUTHORIZED_STATUS = 1
const EXPIRE_STATUS = 2
let authorizationEventSource = null;
const canUseClipboard = navigator.clipboard && window.isSecureContext
$(document).ready(function () {
    hideLoadingView()
    subscribeAuthorizationStatus()
    updateTime()
    $("#authorizationTag").hide()
    $("#endTime").hide()
    $("#deleteAuthorizationBtn").hide()
    setOperation(false)
    setInterval(updateTime, 1000);
    console.log("DOM 加载完成，脚本已运行。");
})

function updateTime() {
    const now = new Date();
    $("#nowTime").text(`当前时间: ${formatTimestamp(now.getTime())}`);
}

function getAuthorizationInfo() {
    $.ajax({
        url: "/api/auth",
        type: "GET",
        success: function (res) {
            if (res.data) {
                $("#authorizationCode").val(res.data)
                $("#deleteAuthorizationBtn").show()
            }
        },
    })
}

function getUniqueCode() {
    $("#getUniqueCodeBtnLoading").show();
    $.ajax({
        url: "/api/uniqueCode",
        type: "GET",
        success: function (res) {
            $("#copyUniqueCodeBtn").prop("disabled", false);
            $("#uniqueCodeText").val(res.data)
        },
        error: function (res) {
            $("#getUniqueCodeBtnLoading").hide();
            if (res.responseJSON){
                alert(`获取唯一标识失败：${res.responseJSON.msg},请联系厂家进行解决`)
            }else{
                alert(`服务器无法连接，请稍后重试`)
            }
        },
        complete: function () {
            $("#getUniqueCodeBtnLoading").hide();
        }
    })
}

function copyUniqueCode() {
    const uniqueCode = $("#uniqueCodeText").val()
    if (uniqueCode === "") {
        alert("无唯一标识")
        return
    }
    if (!canUseClipboard){
        alert("当前无法使用自动复制功能，可能是不在Https标准下，请手动复制")
        return
    }
    navigator.clipboard.writeText(uniqueCode).then(function () {
        alert("复制成功")
    }, function () {
        alert("复制失败")
    });

}

function activateServer() {
    const authorizationCode = $("#authorizationCode").val()
    if (authorizationCode === "") {
        alert("请填写授权码")
        return
    }
    showLoadingView("正在授权，请稍等...")
    $.ajax({
        url: `/api/auth`,
        type: "POST",
        headers: {
            "AuthCode": authorizationCode
        },
        success: function (res) {
            alert(res.msg)
            $("#deleteAuthorizationBtn").show()
        },
        error: function (res) {
            if (res.responseJSON){
                alert(res.responseJSON.msg)
            }else{
                alert(`服务器无法连接，请稍后重试`)
            }
        },
        complete: function () {
            hideLoadingView()
        }
    })
}

function showLoadingView(tip = "加载中") {
    $("#loadingTip").text(tip)
    $("#loadingView").show()
}

function hideLoadingView() {
    $("#loadingView").hide()
}

function clearActivate() {
    showLoadingView("正在清除授权信息，请稍等...")
    $.ajax({
        url: `/api/auth`,
        type: "DELETE",
        success: function (res) {
            alert(res.msg)
            $("#authorizationCode").val("")
            $("#deleteAuthorizationBtn").hide()
        },
        error: function (res) {
            if (res.responseJSON){
                alert(res.responseJSON.msg)
            }else{
                alert(`服务器无法连接，请稍后重试`)
            }
        },
        complete: function () {
            hideLoadingView()
        }
    })
}

function updateAuthorizationStatus(authorizationStatusInfo) {
    $("#authorizationTag").show()
    if (authorizationStatusInfo.tag === AUTHORIZED_STATUS) {
        $("#authorizationTag").removeClass("bg-danger bg-secondary bg-success").addClass("bg-success").text("已授权")
        $("#endTime").show().text("授权到期时间: " + formatTimestamp(authorizationStatusInfo.endTimestamp))
    }
    if (authorizationStatusInfo.tag === UNAUTHORIZED_STATUS) {
        $("#authorizationTag").removeClass("bg-danger bg-secondary bg-success").addClass("bg-secondary").text("未授权")
        $("#endTime").hide()
    }
    if (authorizationStatusInfo.tag === EXPIRE_STATUS) {
        $("#authorizationTag").removeClass("bg-danger bg-secondary bg-success").addClass("bg-danger").text("授权已过期")
        $("#endTime").show().text("授权到期时间: " + formatTimestamp(authorizationStatusInfo.endTimestamp))
    }
}

function subscribeAuthorizationStatus() {
    if (authorizationEventSource) {
        authorizationEventSource.close();
        console.log("已关闭现有 SSE 连接。");
    }

    authorizationEventSource = new EventSource(`/api/status/subscribe`);

    // 监听 'message' 事件 (默认事件类型)
    authorizationEventSource.onmessage = function (event) {
        console.log("收到 SSE 消息:", event.data);
        const response = JSON.parse(event.data);
        updateAuthorizationStatus(response)
    };

    // 监听连接打开事件
    authorizationEventSource.onopen = function (event) {
        console.log("已打开 SSE 连接。")
        $("#subscribeStatus").removeClass("bg-warning bg-success").addClass("bg-success").text("授权状态变化订阅成功")
        setOperation(true)
        getAuthorizationInfo()
    };

    // 监听连接错误事件 (包括连接失败、断开等)
    authorizationEventSource.onerror = function (event) {
        setOperation(false)
        $("#authorizationTag").hide()
        $("#endTime").hide()
        $("#subscribeStatus").removeClass("bg-warning bg-success").addClass("bg-warning").text("等待授权状态变化订阅")
    };
}

function formatTimestamp(timestamp) {
    const date = new Date(timestamp);
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    const hours = String(date.getHours()).padStart(2, '0');
    const minutes = String(date.getMinutes()).padStart(2, '0');
    const seconds = String(date.getSeconds()).padStart(2, '0');
    return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
}

function setOperation(enable){
    $("#deleteAuthorizationBtn").prop("disabled", !enable);
    $("#getUniqueCodeBtn").prop("disabled", !enable);
    $("#activateBtn").prop("disabled", !enable);
}