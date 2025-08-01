const countdownSeconds = 8;
const emptyView = $("#empty-view")
let allList = []
let newestTime = Date.now()
let nowUpdateRecord = null
$(document).ready(function () {
    hideLoadingView()
    updateList()
})
Handlebars.registerHelper('statusClass', function (endTimestamp) {
    return endTimestamp >= newestTime ? 'bg-success' : 'bg-danger';
});
// 用于根据状态返回对应的文本
Handlebars.registerHelper('statusText', function (endTimestamp) {
    console.log(endTimestamp)
    return endTimestamp >= newestTime ? '已授权-到期时间:' : '授权到期-到期时间';

});
Handlebars.registerHelper('formatTimestamp', function (timestamp) {
    return formatTimestamp(timestamp)
});

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

function renderAuthorizationList() {
    $("#authList").empty()
    $("#authRecordCount").text(`授权数量: ${allList.length}个`)
    if (allList.length === 0) {
        emptyView.show()
        return
    }
    newestTime = Date.now()
    emptyView.hide()
    const appDiv = document.getElementById('authList');
    const source = document.getElementById('authListItem').innerHTML;
    const template = Handlebars.compile(source);
    const html = template(allList.filter(item => item.name.includes($("#searchInput").val())));
    appDiv.innerHTML = html;
}

function updateList() {
    $.ajax({
        url: 'http://localhost:1023/api/authList/all',
        type: 'GET',
        success: function (res) {
            allList = res.data
            renderAuthorizationList();
        },
        error: function (res) {
            if (res.responseJSON) {
                alert(res.responseJSON.msg)
            } else {
                alert("服务器无法连接，请重试")
            }
            renderAuthorizationList([]);
        }
    })
}

function addRecord() {
    const name = $("#newRecordName").val();
    if (!name) {
        alert("授权名称不能为空！");
        return;
    }
    const uniqueCode = $("#newRecordUniqueCode").val();
    if (!uniqueCode) {
        alert("机器码不能为空！");
        return;
    }
    const dateString = $("#newRecordDate").val(); // 获取日期字符串，例如 "2025-07-23"
    let timestamp = null;
    if (!dateString) {
        alert("请选择授权截止日期！");
        return;
    }
    const dateObject = new Date(dateString);
    timestamp = dateObject.getTime();
    $.ajax({
        url: 'http://localhost:1023/api/authList/add',
        type: 'POST',
        data: JSON.stringify({
            name: name,
            uniqueCode: uniqueCode,
            endTimestamp: timestamp
        }),
        contentType: 'application/json',
        success: function (res) {
            alert("添加成功");
            updateList();
        },
        error: function (res) {
            if (res.responseJSON) {
                alert(res.responseJSON.msg)
            } else {
                alert(`服务器无法连接，请稍后重试`)
            }
        }
    })
}

function deleteRecord(idStr) {
    const id = parseInt(idStr)
    const nowRecord = allList.find(item => {
        return item.id === id
    })
    if (!nowRecord) {
        alert("找不到操作的记录，奇怪的问题")
        return;
    }
    const confirmResult = confirm(`确定要删除[${nowRecord.name}]吗？`)
    if (!confirmResult) return
    $.ajax({
        url: `http://localhost:1023/api/authList/${id}`,
        type: 'DELETE',
        success: function (res) {
            alert("删除成功");
            updateList()
        },
        error: function (res) {
            if (res.responseJSON) {
                alert(res.responseJSON.msg)
            } else {
                alert(`服务器无法连接，请稍后重试`)
            }
        }
    })
}

function getAuthCode(idStr, btnDom) {
    showLoadingView("正在生成授权码")
    const id = parseInt(idStr)
    $.ajax({
        url: `http://localhost:1023/api/authList/authCode/${id}`,
        type: 'GET',
        success: function (res) {
            navigator.clipboard.writeText(res.data).then(function () {
                alert("自动复制成功")
            }, function () {
                alert("自动复制失败,请从下手动选择复制,8秒后销毁")
                const cardBodyDom = $(btnDom).closest('.card-body');
                createTempAuthCode(res.data, cardBodyDom)
            });
        },
        error: function (res) {
            if (res.responseJSON) {
                alert(res.responseJSON.msg)
            } else {
                alert(`服务器无法连接，请稍后重试`)
            }
        }, complete: function () {
            hideLoadingView()
        }
    })
}

function updateEndTimestamp(idStr) {
    const id = parseInt(idStr)
    nowUpdateRecord = allList.find(item => {
        return item.id === id
    })
    if (!nowUpdateRecord) {
        alert("找不到操作的记录，奇怪的问题")
        return;
    }
    $("#updateRecordName").val(nowUpdateRecord.name);
    $("#updateRecordUniqueCode").val(nowUpdateRecord.uniqueCode);
    const date = new Date(nowUpdateRecord.endTimestamp);
    if (!isNaN(date.getTime())) {
        const year = date.getFullYear();
        const month = String(date.getMonth() + 1).padStart(2, '0');
        const day = String(date.getDate()).padStart(2, '0');
        const formattedDate = `${year}-${month}-${day}`;
        $("#updateRecordDate").val(formattedDate);
    } else {
        $("#updateRecordDate").val('');
    }
    const modal = new bootstrap.Modal($("#updateAuthEndTimeDialog")[0]);
    modal.show();
}

function confirmUpdateRecord() {
    const dateString = $("#updateRecordDate").val(); // 获取日期字符串，例如 "2025-07-23"
    let timestamp = null;
    if (!dateString) {
        alert("请选择授权截止日期！");
        return;
    }
    const newEndTimestamp = new Date(dateString).getTime();
    if (nowUpdateRecord.endTimestamp === newEndTimestamp) {
        alert("新旧日期相同")
        return;
    }
    $.ajax({
        url: 'http://localhost:1023/api/authList/updateEndTimestamp',
        type: 'POST',
        data: JSON.stringify({
            id: nowUpdateRecord.id,
            endTimestamp: newEndTimestamp
        }),
        contentType: 'application/json',
        success: function (res) {
            alert("更新成功");
            updateList()
        },
        error: function (res) {
            if (res.responseJSON) {
                alert(res.responseJSON.msg)
            } else {
                alert("服务器无法连接，请稍后重")
            }
        }
    })
}

function createTempAuthCode(authCode, cardBodyDom) {
    let remainingSeconds = countdownSeconds;
    const tempAuthCodeBox = $('<div>')
        .addClass('temp-auth-code-box');
    const countDownSpan = $('<span>')
        .addClass('count-down')
        .text(countdownSeconds);
    const authCodeSpan = $('<span>')
        .addClass('temp-auth-code')
        .text(`${authCode}`);
    tempAuthCodeBox.append(countDownSpan).append(authCodeSpan);
    cardBodyDom.append(tempAuthCodeBox);
    const countdownInterval = setInterval(() => {
        remainingSeconds--;
        countDownSpan.text(remainingSeconds);
        if (remainingSeconds <= 0) {
            clearInterval(countdownInterval);
            tempAuthCodeBox.remove();
        }
    }, 1000);
}

function showLoadingView(tip = "加载中") {
    $("#loadingTip").text(tip)
    $("#loadingView").show()
}

function clearSearch() {
    $("#searchInput").val("")
    renderAuthorizationList()
}

function searchConfirm() {
    renderAuthorizationList()
}

function hideLoadingView() {
    $("#loadingView").hide()
}


function generateKeys() {
    const modal = new bootstrap.Modal($("#keysGeneratorDialog")[0]);
    modal.show();
    updateKeys()
}

function updateKeys() {
    $.ajax({
        url: '/api/generateKeys',
        type: 'Get',
        success: function (res) {
            const {data} = res
            $("#rsaPublicKey").val(data.rsaPublicKey)
            $("#rsaPrivateKey").val(data.rsaPrivateKey)
            $("#aesKey").val(data.aesKey)
        },
        error: function (res) {
            if (res.responseJSON) {
                alert(res.responseJSON.msg)
            } else {
                alert("服务器无法连接，请稍后重")
            }
        }
    })
}

function copyKeys(id, title) {
    navigator.clipboard.writeText($(`#${id}`).val()).then(function () {
            alert(`${title}自动复制成功`)
        }, function () {
            alert(`${title}自动复制失败,请手动选择复制`)
        }
    )
}